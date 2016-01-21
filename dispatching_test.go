package v8dispatcher

import (
	"fmt"
	"testing"
	"time"
)

type someApi struct{}

func (s someApi) Definition() (string, string, error) {
	return "someApi", `
		someApi = {};
		someApi.now = function() {
			return $sendSync(JSON.stringify({
				"receiver":"someApi",
				"selector":"now"
			}));
		};		
	`, nil
}

func (s someApi) Perform(msg MessageSend) (interface{}, error) {
	if msg.Selector == "now" {
		return time.Now(), nil
	}
	return nil, fmt.Errorf(ErrNoSuchMethod, "someApi", msg.Selector)
}

func TestRequestNow(t *testing.T) {
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
