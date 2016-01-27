package v8dispatcher

import "fmt"

func ExampleMessageDispatcher_CallReturn() {
	md := NewMessageDispatcher()
	md.Worker().Load("ex.js", `
		function now() {
			return new Date();
		}`)
	now, _ := md.CallReturn("this", "now")
	fmt.Println(now)
}
