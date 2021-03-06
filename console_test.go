package v8dispatcher

import "testing"

// clear && go test -v -test.run=TestConsole
func TestConsole(t *testing.T) {
	dist := NewMessageDispatcher()
	capture := &recorder{}
	dist.Register("console", capture)
	err := dist.Worker().Load("console.js", `
		console.log("size",42);
	`)
	if err != nil {
		t.Fatal(err)
	}
	if capture.msg == nil {
		t.Fatal("message not captured")
	}
	if got, want := capture.msg.Selector, "log"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := capture.msg.Arguments[0], "size"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := capture.msg.Arguments[1], float64(42); got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
