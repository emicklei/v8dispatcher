package v8dispatcher

import "testing"

func TestFunctionScheduler(t *testing.T) {
	worker, dist := newWorkerAndDispatcher(t)
	dist.Register("scheduler", NewFunctionScheduler(dist))
	if err := worker.Load("TestFunctionScheduler.js", `
		var utils = {};
		utils.schedule = function(when,then) {
			go_dispatch(function_registry.none, "scheduler", "schedule", when, function_registry.put(then));
		}
		utils.schedule(0,function() {
			$print("performed");
		});
	`); err != nil {
		t.Fatal(err)
	}
}
