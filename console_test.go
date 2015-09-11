package v8dispatcher

import "testing"

func TestConsole(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	capture := &recorder{moduleName: "console"}
	dist.Register(capture)
	err := worker.Load("console.js", `
		console.log("size",42);
	`)
	if err != nil {
		t.Fatal(err)
	}
	if capture.msg == nil {
		t.Fatal("message not captured")
	}
	if got, want := capture.msg.Method, "log"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := capture.msg.Arguments[0], "size"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := capture.msg.Arguments[1], float64(42); got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
