package v8dispatcher

import (
	"fmt"
	"testing"
	"time"
)

func TestCallGoFromJSNoArgsNoReturn(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	rec := &recorder{moduleName: "recorder"}
	dist.Register(rec)
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
	rec := &recorder{moduleName: "recorder"}
	dist.Register(rec)
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
	rec := &recorder{moduleName: "console"}
	dist.Register(rec)
	if err := worker.Load("TestCallJavascriptFromGoNoReturn.js", `
		function calledFromGo(what) {
			console.log(what);
		}
	`); err != nil {
		t.Fatal(err)
	}
	dist.Call("this", "calledFromGo", "hello")
	expectConsoleLogArgument(t, rec, "hello")
}

type someApi struct{}

func (s someApi) ModuleDefinition() (string, string) {
	return "someApi", `
		someApi = {};
		someApi.doit = function(what) {
			go_dispatch(
				function_registry.none,
				"someApi",
				"error",
				what);
		};
		someApi.now = function() {
			return $request(JSON.stringify({
				"receiver":"someApi",
				"method":"now"
			}));
		};		
	`
}

func (s someApi) Perform(msg AsyncMessage) (interface{}, error) {
	if msg.Method == "error" {
		fmt.Printf("go: error was performed with %v\n", msg.Arguments[0])
		return nil, fmt.Errorf("error was performed with %v", msg.Arguments[0])
	}
	if msg.Method == "now" {
		return time.Now(), nil
	}
	return nil, ErrNoSuchMethod
}

func (s someApi) Request(msg MessageSend) (interface{}, error) {
	if msg.Method == "now" {
		return time.Now(), nil
	}
	return nil, ErrNoSuchMethod
}

func TestCallGoInError(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	dist.Register(someApi{})
	if err := worker.Load("TestCallGoInError.js", `
		someApi.doit(42)
	`); err != nil {
		t.Fatal(err)
	}
}

func TestRequestNow(t *testing.T) {
	//t.Skip()
	worker, dist := newWorkerAndDispatcher(t)
	rec := &recorder{moduleName: "console"}
	dist.Register(rec)
	dist.Register(someApi{})
	if err := worker.Load("TestRequestNow.js", `
		console.log(someApi.now())
	`); err != nil {
		t.Fatal(err)
	}
	if len(rec.msg.Arguments[0].(string)) == 0 {
		t.Fail()
	}
}
