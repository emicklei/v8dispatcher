package v8dispatcher

import (
	"testing"
	"time"
)

type someApi struct{}

func (s someApi) ModuleDefinition() (string, string) {
	return "someApi", `
		someApi = {};
		someApi.now = function() {
			return $request(JSON.stringify({
				"receiver":"someApi",
				"method":"now"
			}));
		};		
	`
}

func (s someApi) Perform(msg MessageSend) (interface{}, error) {
	if msg.Selector == "now" {
		return time.Now(), nil
	}
	return nil, ErrNoSuchMethod
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
