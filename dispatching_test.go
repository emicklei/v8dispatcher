package v8dispatcher

import "testing"

func TestCallGoFromJSNoArgsNoReturn(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	rec := new(recorder)
	dist.Register("recorder", rec)
	if err := worker.Load("TestCallGoFromJSNoArgsNoReturn.js", `
		go_dispatch(function_registry.none,"recorder","noargs");
	`); err != nil {
		t.Fatal(err)
	}
	if rec.msg == nil {
		t.Fatal("message not captured")
	}
	if got, want := rec.msg.Method, "noargs"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCallGoFromJSOneArgsNoReturn(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	rec := new(recorder)
	dist.Register("recorder", rec)
	if err := worker.Load("TestCallGoFromJSOneArgsNoReturn.js", `
		go_dispatch(function_registry.none,"recorder","onearg",42);
	`); err != nil {
		t.Fatal(err)
	}
	if rec.msg == nil {
		t.Fatal("message not captured")
	}
	if got, want := rec.msg.Method, "onearg"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := len(rec.msg.Arguments), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := rec.msg.Arguments[0], float64(42); got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestCallJavascriptFromGoNoReturn(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	rec := new(recorder)
	dist.Register("console", rec)
	if err := worker.Load("TestCallJavascriptFromGoNoReturn.js", `
		function calledFromGo(what) {
			console.log(what);
		}
	`); err != nil {
		t.Fatal(err)
	}
	dist.Call("this", "calledFromGo", "hello")
	if rec.msg == nil {
		t.Fatal("message not captured")
	}
	if got, want := rec.msg.Method, "log"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := len(rec.msg.Arguments), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := rec.msg.Arguments[0], "hello"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
