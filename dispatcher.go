package v8dispatcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/emicklei/v8worker"
)

var (
	ErrNoSuchMethod = "%s does not understand %s"
)

// MessageSendHandlerFunc is a function that can be called by the dispatcher if registered using the message selector or receiver.selector.
type MessageSendHandlerFunc func(MessageSend) (interface{}, error)

// MessageSendHandler can be called by the dispatcher if registered using the message receiver.
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

// NewMessageDispatcher returns a new MessageDispatcher initialize with empty handlers and a v8worker.
func NewMessageDispatcher() *MessageDispatcher {
	d := &MessageDispatcher{
		messageHandlerFuncs: map[string]MessageSendHandlerFunc{},
		messageHandlers:     map[string]MessageSendHandler{},
	}
	w := v8worker.New(d.Receive, d.ReceiveSync)
	d.worker = w
	// load scripts
	for _, each := range []string{"js/registry.js", "js/setup.js", "js/console.js"} {
		data, _ := Asset(each)
		if err := w.Load(each, string(data)); err != nil {
			Log("error", "script load error", "source", each, "err", err)
		}
	}
	return d
}

// Worker returns the worker for this dispatcher
func (d *MessageDispatcher) Worker() *v8worker.Worker {
	return d.worker
}

// RegisterFunc adds a function as the handler of a MessageSend.
// The function is called if the name matches the selector of receiver.selector combination.
func (d *MessageDispatcher) RegisterFunc(name string, handler MessageSendHandlerFunc) {
	d.messageHandlerFuncs[name] = handler
}

// Register add a MessageSendHandler implementation that can perform MessageSends.
// The handler is called if the name matches the receiver of the MessageSend.
func (d *MessageDispatcher) Register(name string, handler MessageSendHandler) {
	d.messageHandlers[name] = handler
}

// Send is an asynchronous call to Javascript and does no expect a return value
func (d *MessageDispatcher) Send(receiver string, method string, arguments ...interface{}) error {
	_, err := d.send(MessageSend{
		Receiver:       receiver,
		Selector:       method,
		Arguments:      arguments,
		IsAsynchronous: true,
	})
	return err
}

// SendSync is synchronous call to Javascript and expects a return value
func (d *MessageDispatcher) SendSync(receiver string, method string, arguments ...interface{}) (interface{}, error) {
	return d.send(MessageSend{
		Receiver:       receiver,
		Selector:       method,
		Arguments:      arguments,
		IsAsynchronous: false,
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

// dispatch finds the Go handler registered, calls it and returns the JSON representation of the return value.
// lookup by "receiver" first then "selector" then "receiver.selector" of the message argument.
func (d *MessageDispatcher) dispatch(msg MessageSend) string {
	var result interface{}
	var err error
	if len(msg.Receiver) == 0 {
		performerFunc, ok := d.messageHandlerFuncs[msg.Selector]
		if !ok {
			Log("warn", "no handler func", "selector", msg.Selector)
			return "null"
		}
		result, err = performerFunc(msg)
	} else {
		performer, ok := d.messageHandlers[msg.Receiver]
		if !ok {
			// retry with receiver.selector
			performerFunc, ok := d.messageHandlerFuncs[fmt.Sprintf("%s.%s", msg.Receiver, msg.Selector)]
			if !ok {
				Log("warn", "no handler", "receiver", msg.Receiver, "selector", msg.Selector)
				return "null"
			}
			result, err = performerFunc(msg)
		} else {
			result, err = performer.Perform(msg)
		}
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

func (d *MessageDispatcher) send(ms MessageSend) (interface{}, error) {
	callbackJSON, err := ms.JSON()
	if err != nil {
		Log("error", "message encode failure", "receiver", ms.Receiver, "method", ms.Selector, "err", err)
		return nil, err
	}
	if ms.IsAsynchronous {
		if err := d.worker.Send(callbackJSON); err != nil {
			Log("error", "work send failure", "receiver", ms.Receiver, "method", ms.Selector, "err", err)
			return nil, err
		}
	} else {
		msg := d.worker.SendSync(callbackJSON)
		var value interface{}
		if err := json.Unmarshal([]byte(msg), &value); err != nil {
			return nil, err
		}
		return value, nil
	}
	return nil, nil
}
