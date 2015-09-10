package v8dispatcher

import (
	"encoding/json"
	"strings"

	"github.com/ry/v8worker"
	"gopkg.in/inconshreveable/log15.v2"
)

type Performer interface {
	Perform(msg MessageSend) (interface{}, error)
}

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

// Register is not threadsafe
func (d *MessageDispatcher) Register(receiver string, p Performer) {
	d.performers[receiver] = p
}

// Call dispatches a function in Javascript
func (d *MessageDispatcher) Call(receiver string, method string, arguments ...interface{}) {
	d.send(MessageSend{
		Receiver:  receiver,
		Method:    method,
		Arguments: arguments,
	})
}

// Dispatch is a v8worker handler.
func (d *MessageDispatcher) Dispatch(jsonFromJS string) {
	var msg MessageSend
	if err := json.NewDecoder(strings.NewReader(jsonFromJS)).Decode(&msg); err != nil {
		d.logger.Error("not a valid MessageSend", "err", err)
		return
	}
	performer, ok := d.performers[msg.Receiver]
	if !ok {
		d.logger.Error("unknown receiver", "receiver", msg.Receiver)
		return
	}
	result, err := performer.Perform(msg)

	// all ok, nothing to return
	if err == nil && len(msg.Callback) == 0 {
		return
	}

	var callback MessageSend

	// perform fail, notify javascript about the error
	if err != nil {
		d.logger.Error("perform failure", "receiver", msg.Receiver, "method", msg.Method, "err", err)
		callback = MessageSend{
			Receiver:  "this",
			Method:    "go_error_on_perform",
			Arguments: []interface{}{err.Error()},
		}
	} else {
		// normal return of a value
		callback = MessageSend{
			Receiver:  "this",
			Method:    "callback_dispatch",
			Arguments: []interface{}{msg.Callback, result}, // first argument of callback_dispatch is the functionRef
		}
	}
	d.send(callback)
}

func (d *MessageDispatcher) send(ms MessageSend) error {
	callbackJSON, err := ms.JSON()
	if err != nil {
		d.logger.Error("message encode failure", "receiver", ms.Receiver, "method", ms.Method, "err", err)
		return err
	}
	if err := d.worker.Send(callbackJSON); err != nil {
		d.logger.Error("work send failure", "receiver", ms.Receiver, "method", ms.Method, "err", err)
	}
	return err
}
