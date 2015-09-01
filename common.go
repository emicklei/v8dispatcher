package v8dispatcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ry/v8worker"
	"gopkg.in/inconshreveable/log15.v2"
)

var NoMessageSend = MessageSend{}

type MessageSend struct {
	Receiver  string        `json:"receiver" `
	Method    string        `json:"method" `
	Arguments []interface{} `json:"args" `
	Callback  string        `json:"callback" `
}

func (m MessageSend) String() string {
	return fmt.Sprintf("%s.%s(%v) => %s", m.Receiver, m.Method, m.Arguments, m.Callback)
}

func (m MessageSend) JSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Performer interface {
	Perform(msg MessageSend) (interface{}, error)
}

type PerformError struct {
	Message MessageSend
	Cause   string
}

func (m PerformError) Error() string { return m.Cause }

type MessageDispatcher struct {
	logger     log15.Logger
	performers map[string]Performer
	worker     *v8worker.Worker
}

func NewMessageDispatcher(aLogger log15.Logger) *MessageDispatcher {
	return &MessageDispatcher{
		logger:     aLogger,
		performers: map[string]Performer{},
	}
}

// Worker sets the required v8worker. This cannot be a constructor arg because a worker is created with a handler, the dispatcher itself.
func (d *MessageDispatcher) Worker(worker *v8worker.Worker) {
	d.worker = worker
}

func (d *MessageDispatcher) Register(receiver string, p Performer) {
	d.performers[receiver] = p
}

// Dispatch is a v8worker handler.
func (d *MessageDispatcher) Dispatch(msg string) {
	var ms MessageSend
	if err := json.NewDecoder(strings.NewReader(msg)).Decode(&ms); err != nil {
		d.logger.Error("not a valid MessageSend", "err", err)
		return
	}
	performer, ok := d.performers[ms.Receiver]
	if !ok {
		d.logger.Error("unknown receiver", "receiver", ms.Receiver)
		return
	}
	result, err := performer.Perform(ms)
	if err != nil {
		d.logger.Error("perform failure", "receiver", ms.Receiver, "method", ms.Method, "err", err)
	}
	if len(ms.Callback) == 0 {
		return
	}
	callback := MessageSend{
		Receiver:  "this",
		Method:    "callback_dispatch",
		Arguments: []interface{}{ms.Callback, result}, // first argument of callback_dispatch is the functionRef
	}
	callbackJSON, err := callback.JSON()
	if err != nil {
		d.logger.Error("callback encode failure", "receiver", callback.Receiver, "method", callback.Method, "err", err)
		return
	}
	if err := d.worker.Send(callbackJSON); err != nil {
		d.logger.Error("work send failure", "receiver", callback.Receiver, "method", callback.Method, "err", err)
	}
}
