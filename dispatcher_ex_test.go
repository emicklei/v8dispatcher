package v8dispatcher

import "fmt"

func ExampleMessageDispatcher_CallReturn() {
	md := NewMessageDispatcher()
	now, _ := md.CallReturn("this", "now")
	md.Worker().Load("ex.js", `
		function now() {
			return new Date();
		}`)
	fmt.Println(now)
}
