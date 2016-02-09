package v8dispatcher

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ry/v8worker"
)

// MessageSendHandlerFunc is a function that can be called by the dispatcher if registered using the message selector or receiver.selector.
type MessageSendHandlerFunc func(MessageSend) (interface{}, error)

// MessageSendHandler can be called by the dispatcher if registered using the message receiver.
type MessageSendHandler interface {
	Perform(MessageSend) (interface{}, error)
}

// MessageDispatcher is responsible for handling messages send from Javascript.
// It will do a receiver lookup and perform the message of the receiver.
// If no receiver given then the lookup is based on the selector to find the registered function.
type MessageDispatcher struct {
	messageHandlerFuncs map[string]MessageSendHandlerFunc
	messageHandlers     map[string]MessageSendHandler
	worker              *v8worker.Worker
	traceEnabled        bool
}

// NewMessageDispatcher returns a new MessageDispatcher initialize with empty handlers and a v8worker.
func NewMessageDispatcher() *MessageDispatcher {
	d := &MessageDispatcher{
		messageHandlerFuncs: map[string]MessageSendHandlerFunc{},
		messageHandlers:     map[string]MessageSendHandler{},
		traceEnabled:        false,
	}
	w := v8worker.New(d.Receive, d.ReceiveSync)
	d.worker = w
	// load scripts
	for _, each := range []struct {
		name   string
		source string
	}{
		{"registry.js", registry_js()},
		{"setup.js", setup_js()},
		{"console.js", console_js()},
	} {
		if err := w.Load(each.name, each.source); err != nil {
			Log("error", "script load error", "source", each.name, "err", err)
		}
	}
	// install default console handling
	d.RegisterFunc("console.log", ConsoleLog)
	return d
}

// Worker returns the worker for this dispatcher
func (d *MessageDispatcher) Worker() *v8worker.Worker {
	return d.worker
}

// Trace will cause the internal message sends to be logged. See Log variable.
func (d *MessageDispatcher) Trace(doTrace bool) {
	d.traceEnabled = doTrace
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

// Call is an asynchronous call to Javascript and does no expect a return value
func (d *MessageDispatcher) Call(receiver string, method string, arguments ...interface{}) error {
	_, err := d.send(MessageSend{
		Receiver:       receiver,
		Selector:       method,
		Arguments:      arguments,
		IsAsynchronous: true,
	})
	return err
}

// CallReturn is synchronous call to Javascript and expects a return value
func (d *MessageDispatcher) CallReturn(receiver string, method string, arguments ...interface{}) (interface{}, error) {
	return d.send(MessageSend{
		Receiver:       receiver,
		Selector:       method,
		Arguments:      arguments,
		IsAsynchronous: false,
	})
}

// Callback is an asynchronous call to Javascript that will perform a registered function with optional arguments.
// The funtionReference must have been created with "V8D.function_registry.put(yourFunction)".
func (d *MessageDispatcher) Callback(functionReference string, arguments ...interface{}) error {
	_, err := d.send(MessageSend{
		Receiver:       "V8D",
		Selector:       "callDispatch",
		Arguments:      append([]interface{}{functionReference}, arguments...),
		IsAsynchronous: true,
	})
	return err
}

// Set will add/replace the value for a global variable in Javascript.
func (d *MessageDispatcher) Set(variableName string, value interface{}) error {
	_, err := d.send(MessageSend{
		Receiver:       "V8D",
		Selector:       "set",
		Arguments:      []interface{}{variableName, value},
		IsAsynchronous: true,
	})
	return err
}

// Get will return the value for the global variable in Javascript.
func (d *MessageDispatcher) Get(variableName string) (interface{}, error) {
	return d.send(MessageSend{
		Receiver:       "V8D",
		Selector:       "get",
		Arguments:      []interface{}{variableName},
		IsAsynchronous: false,
	})
}

// ReceiveSync is a v8worker send sync handler.
func (d *MessageDispatcher) ReceiveSync(jsonFromJS string) string {
	if d.traceEnabled {
		Log("trace", "ReceiveSync", "json", jsonFromJS)
	}
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
	if d.traceEnabled {
		Log("trace", "Receive", "json", jsonFromJS)
	}
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
	if d.traceEnabled {
		Log("trace", "dispatch", "msg", msg)
	}
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
		Log("error", "perform failed", "err", err.Error())
		return err.Error() // TODO
	}

	// if no return value is expected and no callback is requested then we are done
	if msg.IsAsynchronous && len(msg.Callback) == 0 {
		return ""
	}

	// make the JSON for the result
	data, err := json.Marshal(result)
	if err != nil {
		Log("error", "marshal error", "err", err.Error())
		return err.Error() // TODO
	}

	// if a callback is given then call this first with the result
	if len(msg.Callback) > 0 {
		callDispatch := MessageSend{
			Receiver:       "V8D",
			Selector:       "callDispatch",
			Arguments:      append([]interface{}{msg.Callback}, string(data)),
			IsAsynchronous: true,
		}
		_, err := d.send(callDispatch)
		if err != nil {
			Log("error", "callDispatch failed", "err", err.Error())
			return err.Error() // TODO
		}
	}
	return string(data)
}

// send will perform a MessageSend in Javascript
// if the message is synchronous then return the result of the Javascript function.
func (d *MessageDispatcher) send(msg MessageSend) (interface{}, error) {
	if d.traceEnabled {
		Log("trace", "send", "msg", msg)
	}
	callbackJSON, err := msg.JSON()
	if err != nil {
		Log("error", "message encode failure", "receiver", msg.Receiver, "method", msg.Selector, "err", err)
		return nil, err
	}
	if msg.IsAsynchronous {
		if err := d.worker.Send(callbackJSON); err != nil {
			Log("error", "work send failure", "receiver", msg.Receiver, "method", msg.Selector, "err", err)
			return nil, err
		}
		return nil, nil
	}
	// synchronous
	reply := d.worker.SendSync(callbackJSON)
	var value interface{}
	if err := json.Unmarshal([]byte(reply), &value); err != nil {
		Log("error", "unmarshal Javascript message failure", "msg", msg, "reply", reply, "err", err)
		return nil, err
	}
	return value, nil
}
