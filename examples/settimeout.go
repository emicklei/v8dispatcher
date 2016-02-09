/*
This example shows how you can use callbacks created in Javascript and called from Go.
From the Window setTimeout documentation:

	The setTimeout() method calls a function or evaluates an expression after a specified number of milliseconds.

Using a MessageDispatcher, a function is registered for the selector "setTimeout".
When called from Javascript, it starts a go-routine with a sleep and the callback invocation.
To make this function available in Javascript, a definition is loaded which calls the V8D put function.

*/
package main

import (
	"time"

	v8d "github.com/emicklei/v8dispatcher"
)

func main() {
	m := v8d.NewMessageDispatcher()
	m.Trace(true)
	m.RegisterFunc("setTimeout", func(msg v8d.MessageSend) (interface{}, error) {
		fnc := msg.Arguments[0].(string)
		ms := time.Duration(msg.Arguments[1].(float64)) * time.Millisecond
		go func() {
			time.Sleep(ms)
			m.Callback(fnc)
		}()
		return nil, nil
	})
	m.Worker().Load("setTimeout.js", `
		function setTimeout(func,ms){
			V8D.call("","setTimeout",V8D.function_registry.put(func),ms);
		}	
	`)
	m.Worker().Load("test.js", `
		setTimeout(function(){
			console.log("timed out");
		}, 1000);
	`)
	time.Sleep(1500 * time.Millisecond)
}
