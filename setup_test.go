package v8dispatcher

import (
	"io/ioutil"
	"testing"

	"github.com/emicklei/v8worker"
)

func newWorkerAndDispatcher(t *testing.T) (*v8worker.Worker, *MessageDispatcher) {
	dist := NewMessageDispatcher()
	for _, each := range []string{"registry.js", "setup.js", "console.js"} {
		//t.Log("reading " + each)
		src, err := ioutil.ReadFile(each)
		if err != nil {
			t.Fatal(err)
		}
		//t.Log("loading " + each)
		err = dist.Worker().Load(each, string(src))
		if err != nil {
			t.Fatal(err)
		}
	}
	return dist.Worker(), dist
}

func benchNewWorkerAndDispatcher(b *testing.B) (*v8worker.Worker, *MessageDispatcher) {
	dist := NewMessageDispatcher()
	for _, each := range []string{"registry.js", "setup.js", "console.js"} {
		//t.Log("reading " + each)
		src, err := ioutil.ReadFile(each)
		if err != nil {
			b.Fatal(err)
		}
		//t.Log("loading " + each)
		err = dist.Worker().Load(each, string(src))
		if err != nil {
			b.Fatal(err)
		}
	}
	return dist.Worker(), dist
}

type recorder struct {
	moduleName string
	source     string
	msg        *MessageSend
}

func (r *recorder) Perform(msg MessageSend) (interface{}, error) {
	r.msg = &msg
	return nil, nil
}

func expectConsoleLogArgument(t *testing.T, rec *recorder, arg interface{}) {
	if rec.msg == nil {
		t.Fatal("message not recorded")
	}
	if got, want := rec.msg.Selector, "log"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := len(rec.msg.Arguments), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := rec.msg.Arguments[0], arg; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
