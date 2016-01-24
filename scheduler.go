package v8dispatcher

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"
)

// scheduledCall encapsulates a Javascript function and its arguments to call it a specified time.
type scheduledCall struct {
	when           time.Time
	message        MessageSend
	previous, next *scheduledCall
}

// FunctionScheduler hold MessageSend values that are performed after a specified time.
type FunctionScheduler struct {
	mutex      *sync.RWMutex
	head, tail *scheduledCall
	dispatcher *MessageDispatcher
}

// NewFunctionScheduler returns a new FunctionScheduler and registers itself as a MessageHandler.
func NewFunctionScheduler(dispatcher *MessageDispatcher) *FunctionScheduler {
	fs := &FunctionScheduler{
		mutex:      new(sync.RWMutex),
		dispatcher: dispatcher,
	}
	dispatcher.Worker().Load("FunctionScheduler.js", fs.source())
	dispatcher.Register("V8D.scheduler", fs)
	return fs
}

// source returns the Javascript initialization to provide the scheduling function.
func (s *FunctionScheduler) source() string {
	return `
		V8D.schedule = function(after,then) {			
			V8D.call("V8D.scheduler", "schedule", after, V8D.function_registry.put(then));
		}
	`
}

// Perform handles a message send from Javascript. FunctionScheduler is a MessageHandler.
func (s *FunctionScheduler) Perform(msg MessageSend) (interface{}, error) {
	if Debug {
		Log("debug", "Perform", "msg", msg)
	}
	if "schedule" != msg.Selector {
		return nil, fmt.Errorf(ErrNoSuchMethod, "V8D.scheduler", msg.Selector)
	}
	if len(msg.Arguments) != 2 {
		return nil, errors.New("expected `after` and `then` arguments")
	}
	when, ok := msg.Arguments[0].(float64)
	if !ok {
		return nil, errors.New("first argument `after` must be delay in milliseconds (number)")
	}
	then, ok := msg.Arguments[1].(string)
	if !ok {
		return nil, errors.New("second argument `then` must be a function reference (string)")
	}
	scheduledMsg := MessageSend{
		Receiver:       "V8D",
		Selector:       "callDispatch",
		Arguments:      []interface{}{then},
		IsAsynchronous: true,
	}
	return nil, s.Schedule(int64(when), scheduledMsg)
}

// Reset forgets about all scheduled calls.
func (s *FunctionScheduler) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.head = nil
	s.tail = nil
}

// PerformCallsBefore performs all calls scheduled before a certain point in time.
// Each call is run in its own goroutine
func (s *FunctionScheduler) PerformCallsBefore(when time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for s.head != nil && when.After(s.head.when) {
		// first detach then perform to allow new call to be inserted
		call := s.head
		s.head = call.next
		go s.dispatcher.send(call.message)
	}
}

// Schedule adds a call to be performed in the future.
func (s *FunctionScheduler) Schedule(delayInMilliseconds int64, msg MessageSend) error {
	if delayInMilliseconds < 0 {
		return errors.New("cannot schedule a function call in the past")
	}
	if delayInMilliseconds == 0 {
		go s.dispatcher.send(msg)
		return nil
	}
	absoluteTime := time.Now().Add(time.Duration(delayInMilliseconds) * time.Millisecond)
	call := &scheduledCall{
		when:    absoluteTime,
		message: msg,
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.head == nil {
		s.head = call
		s.tail = call
		return nil
	}
	if s.head.when.After(call.when) {
		// new head
		s.head.previous = call
		call.next = s.head
		s.head = call
		return nil
	}
	if call.when.After(s.tail.when) {
		// new tail
		call.previous = s.tail
		s.tail.next = call
		s.tail = call
		return nil
	}
	// on or between head and tail
	here := s.head.next
	for call.when.After(here.when) {
		here = here.next
	}
	// here is after call, it must be scheduled before it
	call.previous = here.previous
	call.next = here
	here.previous.next = call
	here.previous = call
	return nil
}

// String is for debugging
func (s *FunctionScheduler) String() string {
	var buf bytes.Buffer
	here := s.head
	for here != nil {
		buf.WriteString(fmt.Sprintf("\n -> %v", here.when))
		here = here.next
	}
	return buf.String()
}
