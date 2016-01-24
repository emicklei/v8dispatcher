package v8dispatcher

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func setupLog(t *testing.T) {
	// COPIED from std Log definition
	Log = func(level, msg string, kvs ...interface{}) {
		// default uses standard logging
		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, "[%s] %s", level, msg)
		for i := 0; i < len(kvs); i = i + 2 {
			var v interface{}
			if len(kvs) == i+1 {
				v = "*** missing ***"
			} else {
				v = kvs[i+1]
			}
			fmt.Fprintf(buf, ", %v = %v", kvs[i], v)
		}
		t.Log(buf.String())
	}
}

func TestFunctionSchedulerImmediate(t *testing.T) {
	setupLog(t)
	worker, dist := newWorkerAndDispatcher(t)
	rec := &recorder{}
	_ = NewFunctionScheduler(dist)
	dist.Register("console", rec)
	if err := worker.Load("TestFunctionScheduler.js", `		
		V8D.schedule(0,function() {
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
	rec := &recorder{}
	dist.Register("console", rec)
	if err := worker.Load("TestFunctionScheduler.js", `		
		V8D.schedule(100,function() {
			console.log("performed 100 ms later");
		});
	`); err != nil {
		t.Fatal(err)
	}
	s.PerformCallsBefore(time.Now().Add(1 * time.Second))
	time.Sleep(1 * time.Second)
	expectConsoleLogArgument(t, rec, "performed 100 ms later")
}
