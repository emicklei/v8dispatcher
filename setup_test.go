package v8dispatcher

import (
	"io/ioutil"
	"testing"

	"bitbucket.org/emicklei/firespark/logger"

	"gopkg.in/inconshreveable/log15.v2"

	"github.com/ry/v8worker"
)

func TestConsole(t *testing.T) {
	dist := NewMessageDispatcher(log15.New())
	worker := v8worker.New(dist.Dispatch)
	dist.Worker(worker)

	dist.Register("console", Console{})
	dist.Register("echo", Echo{})

	src, err := ioutil.ReadFile("setup.js")
	if err != nil {
		t.Fatal(err)
	}
	err = worker.Load("setup.js", string(src))
	if err != nil {
		t.Fatal(err)
	}

	err = worker.Load("console.js", `
		console.log("size",42);
		
		function putit_togo(arg) {
			go_dispatch(undefined, "echo", "return", arg);
		}		
		
		function getit_fromgo(then) {
			go_dispatch(function_registry.put(then), "echo", "return", 42);
		}
		getit_fromgo(function(msg){
			$print(msg);
		});		
		
		putit_togo(36)
	`)
	if err != nil {
		t.Fatal(err)
	}
}

type Echo struct{}

func (e Echo) Perform(msg MessageSend) (interface{}, error) {
	logger.Logger.Info("perform", "msg", msg.String())
	return 21, nil
}
