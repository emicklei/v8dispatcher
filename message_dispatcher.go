package v8dispatcher

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/emicklei/v8worker"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	ErrNoSuchMethod = errors.New("no such method")
)

// Module represents a Javascript object with functions that call into its Go counterpart.
type Module interface {
	// ModuleDefinition returns the name of the module as it will be known in Javascript
	// and Javascript source to create this module (global variable).
	// It returns an error if loading the source failed.
	ModuleDefinition() (string, string)

	// Perform will call the function associated to the Method of the message.
	// It returns a value (optionally) and an error if the call failed.
	Perform(msg AsyncMessage) (interface{}, error)
}

// MessageDispatcher is responsible for handling messages send from Javascript.
// It will do a receiver lookup and perform the messages by the receiver.
type MessageDispatcher struct {
	logger     log15.Logger
	performers map[string]Module
	worker     *v8worker.Worker
}

func NewMessageDispatcher(aLogger log15.Logger) *MessageDispatcher {
	return &MessageDispatcher{
		logger:     aLogger,
		performers: map[string]Module{},
	}
}

// Worker sets the required v8worker. This cannot be a constructor arg because a worker is created with a handler, the dispatcher itself.
func (d *MessageDispatcher) Worker(worker *v8worker.Worker) {
	d.worker = worker
}

// Register adds a Module and makes it available to Javascript by its defintion name.
// Not yet threadsafe
func (d *MessageDispatcher) Register(p Module) error {
	name, source := p.ModuleDefinition()
	if len(source) > 0 {
		if err := d.worker.Load("v8dispatcher_"+name+".js", source); err != nil {
			d.logger.Error("module load failed", "module", name, "err", err.Error())
			return err
		}
	}
	d.performers[name] = p
	return nil
}

// Call dispatches a function in Javascript
func (d *MessageDispatcher) Call(receiver string, method string, arguments ...interface{}) {
	d.send(MessageSend{
		Receiver:  receiver,
		Method:    method,
		Arguments: arguments,
	})
}

// DispatchRequest is a v8worker exchange handler.
func (d *MessageDispatcher) DispatchRequest(jsonFromJS string) string {
	return "42"
}

// DispatchSend is a v8worker callback handler.
func (d *MessageDispatcher) DispatchSend(jsonFromJS string) {
	var msg AsyncMessage
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
	var callback MessageSend
	if err != nil {
		d.logger.Error(err.Error())
		for ix, each := range strings.Split(msg.Stack, "\n") {
			if ix != 1 {
				d.logger.Error(each)
			}
		}
		return
	} else {
		// check onReturn
		if len(msg.Callback) == 0 {
			if result != nil {
				d.logger.Error("perform returned result but no callback was given", "receiver",
					msg.Receiver, "method", msg.Method, "result", result)
				return
			}
			return
		}
		callback = MessageSend{
			Receiver:  "this",
			Method:    "callback_dispatch",
			Arguments: []interface{}{msg.Callback, result}, // first argument of callback_dispatch is the functionRef
		}
	}
	d.send(callback)
}

func (d *MessageDispatcher) send(ms MessageSend) {
	callbackJSON, err := ms.JSON()
	if err != nil {
		d.logger.Error("message encode failure", "receiver", ms.Receiver, "method", ms.Method, "err", err)
		return
	}
	if err := d.worker.Send(callbackJSON); err != nil {
		d.logger.Error("work send failure", "receiver", ms.Receiver, "method", ms.Method, "err", err)
	}
}
