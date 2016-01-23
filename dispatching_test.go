package v8dispatcher

import (
	"testing"
	"time"
)

func TestRequestNow(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	rec := &recorder{}
	dist.Register("console", rec)

	worker.Load("someApi.js", `
		someApi = {};
		someApi.now = function() {
			return V8D.callReturn("","someApi.now");
		};		
	`)

	dist.RegisterFunc("someApi.now", func(msg MessageSend) (interface{}, error) {
		return time.Now(), nil
	})

	if err := worker.Load("TestRequestNow.js", `
		console.log(someApi.now())
	`); err != nil {
		t.Fatal(err)
	}
	if rec.msg == nil {
		t.Fatal("no msg recorded")
	}
	if len(rec.msg.Arguments) == 0 {
		t.Fatal("no arguments recorded")
	}
	s, ok := rec.msg.Arguments[0].(string)
	if !ok {
		t.Fatal("string expected")
	}
	if len(s) == 0 {
		t.Fail()
	}
	t.Logf("%#v", rec.msg)
}

func BenchmarkRequestFromGo(b *testing.B) {
	worker, _ := benchNewWorkerAndDispatcher(b)
	if err := worker.Load("BenchmarkRequestFromGo.js", `
		function dummy(what) {
			return what;
		}
	`); err != nil {
		b.Fatal(err)
	}
	msg := MessageSend{
		Receiver:  "this",
		Selector:  "dummy",
		Arguments: []interface{}{42},
	}
	js, _ := msg.JSON()
	for n := 0; n < b.N; n++ {
		worker.SendSync(js)
	}
}
