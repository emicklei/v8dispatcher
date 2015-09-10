package v8dispatcher

import (
	"io/ioutil"
	"testing"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/ry/v8worker"
)

func newWorkerAndDispatcher(t *testing.T) (*v8worker.Worker, *MessageDispatcher) {
	dist := NewMessageDispatcher(log15.New())
	worker := v8worker.New(dist.Dispatch)
	dist.Worker(worker)
	//t.Log("reading setup.js")
	src, err := ioutil.ReadFile("setup.js")
	if err != nil {
		t.Fatal(err)
	}
	//t.Log("loading setup.js")
	err = worker.Load("setup.js", string(src))
	if err != nil {
		t.Fatal(err)
	}
	return worker, dist
}

type recorder struct {
	msg *MessageSend
}

func (r *recorder) Perform(msg MessageSend) (interface{}, error) {
	r.msg = &msg
	return nil, nil
}

func expectConsoleLogArgument(t *testing.T, rec *recorder, arg interface{}) {
	if rec.msg == nil {
		t.Fatal("message not recorded")
	}
	if got, want := rec.msg.Method, "log"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := len(rec.msg.Arguments), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := rec.msg.Arguments[0], arg; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
