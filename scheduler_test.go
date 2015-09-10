package v8dispatcher

import (
	"testing"
	"time"
)

func TestFunctionSchedulerImmediate(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	rec := new(recorder)
	dist.Register("console", rec)
	dist.Register("scheduler", NewFunctionScheduler(dist))
	if err := worker.Load("TestFunctionScheduler.js", `
		var utils = {};
		
		// define
		utils.schedule = function(when,then) {
			go_dispatch(function_registry.none, "scheduler", "schedule", when, function_registry.put(then));
		}
		
		// use
		utils.schedule(0,function() {
			console.log("performed immediately");
		});
	`); err != nil {
		t.Fatal(err)
	}
	// because the the console.log (v8->go) is asynchronous, it won't be there really immediately.
	time.Sleep(1 * time.Second)
	expectConsoleLogArgument(t, rec, "performed immediately")
}

func TestFunctionScheduler100ms(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	s := NewFunctionScheduler(dist)
	dist.Register("scheduler", s)
	rec := new(recorder)
	dist.Register("console", rec)
	if err := worker.Load("TestFunctionScheduler.js", `
		var utils = {};
		utils.schedule = function(when,then) {
			go_dispatch(function_registry.none, "scheduler", "schedule", when, function_registry.put(then));
		}
		utils.schedule(100,function() {
			console.log("performed 100 ms later");
		});
	`); err != nil {
		t.Fatal(err)
	}
	s.PerformCallsBefore(time.Now().Add(1 * time.Second))
	time.Sleep(1 * time.Second)
	expectConsoleLogArgument(t, rec, "performed 100 ms later")
}
