package v8dispatcher

import (
	"testing"
	"time"
)

var someApiSrc = `
		someApi = {};
		
		someApi.now = function() {
			return V8D.callReturn("someApi","now");
		};	
		
		someApi.callThen = function() {
			V8D.callThen("someApi","now", function() {
				console.log("callThen");
			});
		};	
		
		someApi.callThenArgument = function() {
			V8D.callThen("someApi","now", function(argument) {
				console.log("callThen with",argument);
			});
		};			
	`

func TestCallReturn(t *testing.T) {
	dist := NewMessageDispatcher()
	rec := &recorder{}
	dist.Register("console", rec)
	dist.Worker().Load("someApi.js", someApiSrc)

	dist.RegisterFunc("someApi.now", func(msg MessageSend) (interface{}, error) {
		return time.Now(), nil
	})

	if err := dist.Worker().Load("TestRequestNow.js", `
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

func TestCallThen(t *testing.T) {
	dist := NewMessageDispatcher()
	rec := &recorder{}
	dist.Register("console", rec)
	dist.Worker().Load("someApi.js", someApiSrc)

	dist.RegisterFunc("someApi.now", func(msg MessageSend) (interface{}, error) {
		return time.Now(), nil
	})

	if err := dist.Worker().Load("TestRequestNow.js", `
		someApi.callThen()
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
	if s != "callThen" {
		t.Fail()
	}
	t.Logf("%#v", rec.msg)
}

func BenchmarkRequestFromGo(b *testing.B) {
	dist := NewMessageDispatcher()
	worker := dist.Worker()
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
