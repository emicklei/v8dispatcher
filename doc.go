/*
Package v8dispatcher adds a message layer on top of the v8worker package (binding to Javascript V8).

The v8worker package is a simple binding that provides a few Javascript operations to send simple messages (string) to and receive from Go.
Recently, the v8worker package has been enhanced to support synchronous communication; this allows for accessing return values from functions.
The v8dispatcher package sends MessageSend values serialized as JSON strings to be dispatched in Go or Javascript.
A MessageDispatcher is used to dispatch MessageSend values to function calls, both in Go and in Javascript (no reflection).

Methods available in Go to invoke custom functions in Javascript (see MessageDispatcher):

	// Call is an asynchronous call to Javascript and does no expect a return value
	// CallReturn is synchronous call to Javascript and expects a return value

Functions available in Javascript to invoke custom functions in Go (see js folder):

	// V8D.call performs a MessageSend in Go and does NOT return a value.
	// V8D.callReturn performs a MessageSend in Go and returns the value from that result
	// V8D.callThen performs a MessageSend in Go which can call the onReturn function.

Variables in Javascript can be set and get using:

	// Set will add/replace the value for a global variable in Javascript.
	// Get will return the value for the global variable in Javascript.

Dispatching MessageSend values to functions in Go requires the registration of handlers.

The `MessageSendHandlerFunc` can be used to map a function name (the MessageSend selector) to a Go function.
Use the method ` RegisterFunc(name string, handler MessageSendHandlerFunc)` on the dispatcher.
By implementing the `MessageHandler` interface, the mapping from selectors is implemented in the Perform method.
Use the method `Register(name string, handler MessageSendHandler)` that register a handler.

For examples see the README.md and the tests.

(c) 2016, http://ernestmicklei.com. MIT License
*/
package v8dispatcher
