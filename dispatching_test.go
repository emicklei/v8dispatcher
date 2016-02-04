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
		
		someApi.testCallThen = function() {
			V8D.callThen("someApi","now", function() {
				console.log("callThen");
			});
		};	
		
		someApi.testCallThenArgument = function() {
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
	//t.Logf("%#v", rec.msg)
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
		someApi.testCallThen()
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
	//t.Logf("%#v", rec.msg)
}

func TestCallThenWithArgument(t *testing.T) {
	dist := NewMessageDispatcher()
	rec := &recorder{}
	dist.Register("console", rec)
	dist.Worker().Load("someApi.js", someApiSrc)

	dist.RegisterFunc("someApi.now", func(msg MessageSend) (interface{}, error) {
		return time.Now(), nil
	})

	if err := dist.Worker().Load("TestRequestNow.js", `
		someApi.testCallThenArgument()
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
	if s != "callThen with" {
		t.Errorf("got %s", s)
	}
	t.Logf("argument=%#v", rec.msg.Arguments[1])
}

func TestSetGet(t *testing.T) {
	dist := NewMessageDispatcher()
	rec := &recorder{}
	dist.Register("console", rec)
	dist.Set("SomeVar", 42)
	if err := dist.Worker().Load("TestSet.js", `
		console.log(this["SomeVar"]);
	`); err != nil {
		t.Fatal(err)
	}
	if rec.msg == nil {
		t.Fatal("no msg recorded")
	}
	if len(rec.msg.Arguments) == 0 {
		t.Fatal("no arguments recorded")
	}
	i, ok := rec.msg.Arguments[0].(float64)
	if !ok {
		t.Fatalf("float64 expected, got %T", rec.msg.Arguments[0])
	}
	if i != 42 {
		t.Fail()
	}
	t.Logf("%#v", rec.msg)

	v, err := dist.Get("SomeVar")
	if err != nil {
		t.Fatal(err)
	}
	if i != v {
		t.Fail()
	}
}

func TestRoundTripWithMap(t *testing.T) {
	dist := NewMessageDispatcher()
	var gotArgument interface{}
	var gotReturn interface{}
	dist.RegisterFunc("putBasket", func(msg MessageSend) (interface{}, error) {
		gotArgument = msg.Arguments[0]
		return nil, nil
	})
	if err := dist.Worker().Load("TestRoundTripWithMap.js", `
		function getBasket(basket){
			V8D.call("","putBasket",basket);
			return basket
		};
	`); err != nil {
		t.Fatal(err)
	}
	gotReturn, err := dist.CallReturn("this", "getBasket", map[string]interface{}{"size": 42})
	if err != nil {
		t.Fatal(err)
	}
	if gotArgument == nil {
		t.Error("argument not caught")
	}
	if gotReturn == nil {
		t.Error("return not passes")
	}
	mapA, ok := gotArgument.(map[string]interface{})
	if !ok {
		t.Error("map argument expected")
	}
	mapR, ok := gotReturn.(map[string]interface{})
	if !ok {
		t.Error("map return expected")
	}
	if mapA["size"] != float64(42) {
		t.Errorf("42 expected, got %v", mapA["size"])
	}
	if mapR["size"] != float64(42) {
		t.Errorf("42 expected, got %v", mapR["size"])
	}
	t.Log(gotArgument, gotReturn)
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
