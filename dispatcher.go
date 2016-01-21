package v8dispatcher

import (
	"encoding/json"
	"strings"

	"github.com/emicklei/v8worker"
)

var (
	ErrNoSuchMethod = "%s does not understand %s"
)

type MessageSendHandlerFunc func(MessageSend) (interface{}, error)

type MessageSendHandler interface {
	Perform(MessageSend) (interface{}, error)
}

// MessageDispatcher is responsible for handling messages send from Javascript.
// It will do a receiver lookup and perform the messages by the receiver.
type MessageDispatcher struct {
	messageHandlerFuncs map[string]MessageSendHandlerFunc
	messageHandlers     map[string]MessageSendHandler
	worker              *v8worker.Worker
}

func NewMessageDispatcher() *MessageDispatcher {
	d := &MessageDispatcher{
		messageHandlerFuncs: map[string]MessageSendHandlerFunc{},
		messageHandlers:     map[string]MessageSendHandler{},
	}
	w := v8worker.New(d.Receive, d.ReceiveSync)
	d.worker = w
	return d
}

func (d *MessageDispatcher) Worker() *v8worker.Worker {
	return d.worker
}

func (d *MessageDispatcher) RegisterFunc(name string, handler MessageSendHandlerFunc) {
	d.messageHandlerFuncs[name] = handler
}

func (d *MessageDispatcher) Register(name string, handler MessageSendHandler) {
	d.messageHandlers[name] = handler
}

// Call dispatches a function in Javascript
func (d *MessageDispatcher) Call(receiver string, method string, arguments ...interface{}) {
	d.send(MessageSend{
		Receiver:  receiver,
		Selector:  method,
		Arguments: arguments,
	})
}

// ReceiveSync is a v8worker send sync handler.
func (d *MessageDispatcher) ReceiveSync(jsonFromJS string) string {
	var msg MessageSend
	if err := json.NewDecoder(strings.NewReader(jsonFromJS)).Decode(&msg); err != nil {
		Log("error", "not a valid MessageSend", "err", err)
		return err.Error() // TODO
	}
	msg.IsAsynchronous = false
	return d.dispatch(msg)
}

// Receive is a v8worker send async handler.
func (d *MessageDispatcher) Receive(jsonFromJS string) {
	var msg MessageSend
	if err := json.NewDecoder(strings.NewReader(jsonFromJS)).Decode(&msg); err != nil {
		Log("error", "not a valid MessageSend", "err", err)
		return
	}
	msg.IsAsynchronous = true
	_ = d.dispatch(msg)
}

func (d *MessageDispatcher) dispatch(msg MessageSend) string {
	var result interface{}
	var err error
	if len(msg.Receiver) == 0 {
		performer, ok := d.messageHandlerFuncs[msg.Selector]
		if !ok {
			Log("error", "no handler func", "selector", msg.Selector)
			return "" // TODO
		}
		result, err = performer(msg)
	} else {
		performer, ok := d.messageHandlers[msg.Receiver]
		if !ok {
			Log("error", "no handler", "receiver", msg.Receiver)
			return "" // TODO
		}
		result, err = performer.Perform(msg)
	}
	if err != nil {
		Log("error", err.Error())
		return err.Error() // TODO
	}
	data, err := json.Marshal(result)
	if err != nil {
		Log("error", err.Error())
		return err.Error() // TODO
	}
	return string(data)
}

func (d *MessageDispatcher) send(ms MessageSend) {
	callbackJSON, err := ms.JSON()
	if err != nil {
		Log("error", "message encode failure", "receiver", ms.Receiver, "method", ms.Selector, "err", err)
		return
	}
	if err := d.worker.Send(callbackJSON); err != nil {
		Log("error", "work send failure", "receiver", ms.Receiver, "method", ms.Selector, "err", err)
	}
}
