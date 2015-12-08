package v8dispatcher

import (
	"encoding/json"
	"strings"

	"github.com/emicklei/v8worker"
	"gopkg.in/inconshreveable/log15.v2"
)

var (
	ErrNoSuchMethod = "%s does not understand %s"
)

// MessageDispatcher is responsible for handling messages send from Javascript.
// It will do a receiver lookup and perform the messages by the receiver.
type MessageDispatcher struct {
	logger          log15.Logger
	messageHandlers map[string]Module
	worker          *v8worker.Worker
}

func NewMessageDispatcher(aLogger log15.Logger) *MessageDispatcher {
	return &MessageDispatcher{
		logger:          aLogger,
		messageHandlers: map[string]Module{},
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
	d.messageHandlers[name] = p
	return nil
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
		d.logger.Error("not a valid MessageSend", "err", err)
		return err.Error() // TODO
	}
	msg.IsAsynchronous = false
	return d.dispatch(msg)
}

// Receive is a v8worker send async handler.
func (d *MessageDispatcher) Receive(jsonFromJS string) {
	var msg MessageSend
	if err := json.NewDecoder(strings.NewReader(jsonFromJS)).Decode(&msg); err != nil {
		d.logger.Error("not a valid MessageSend", "err", err)
		return
	}
	msg.IsAsynchronous = true
	_ = d.dispatch(msg)
}

func (d *MessageDispatcher) dispatch(msg MessageSend) string {
	performer, ok := d.messageHandlers[msg.Receiver]
	if !ok {
		d.logger.Error("unknown receiver", "receiver", msg.Receiver)
		return "" // TODO
	}
	result, err := performer.Perform(msg)
	if err != nil {
		d.logger.Error(err.Error())
		return err.Error() // TODO
	}
	data, err := json.Marshal(result)
	if err != nil {
		d.logger.Error(err.Error())
		return err.Error() // TODO
	}
	return string(data)
}

func (d *MessageDispatcher) send(ms MessageSend) {
	callbackJSON, err := ms.JSON()
	if err != nil {
		d.logger.Error("message encode failure", "receiver", ms.Receiver, "method", ms.Selector, "err", err)
		return
	}
	if err := d.worker.Send(callbackJSON); err != nil {
		d.logger.Error("work send failure", "receiver", ms.Receiver, "method", ms.Selector, "err", err)
	}
}
