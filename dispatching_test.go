package v8dispatcher

import "testing"

func TestCallGoFromJSNoArgsNoReturn(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	rec := new(recorder)
	dist.Register("recorder", rec)
	if err := worker.Load("console.js", `
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
	if err := worker.Load("console.js", `
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

//func TestConsole(t *testing.T) {
//	worker, dist := newWorkerAndDispatcher(t)
//	dist.Register("console", Console{})

//	src, err := ioutil.ReadFile("setup.js")

//	err = worker.Load("console.js", `
//		console.log("size",42);

//		function putit_togo(arg) {
//			go_dispatch(function_registry.none, "echo", "noreturn", arg);
//		}

//		function getit_fromgo(then) {
//			go_dispatch(function_registry.put(then), "echo", "return", 42);
//		}
//		getit_fromgo(function(msg){
//			$print(msg);
//		});

//		putit_togo(36)

//		go_dispatch(function_registry.none, "badthings", "happen", "today");
//	`)
//	if err != nil {
//		t.Fatal(err)
//	}
//}
