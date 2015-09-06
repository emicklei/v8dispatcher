package v8dispatcher

import (
	"errors"
	"log"
	"testing"
)

func TestHandleErrorInGo(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	dist.Register("badthings", BadThings{})

	err := worker.Load("TestHandleErrorInGo.js", `
		go_dispatch(function_registry.none, "badthings", "happen", "today");
	`)
	if err != nil {
		t.Fatal(err)
	}
}

type BadThings struct{}

func (b BadThings) Perform(msg MessageSend) (interface{}, error) {
	log.Println("perform", "msg", msg.String())
	return nil, errors.New("something bad happened")
}
