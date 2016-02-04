/*
Package v8dispatcher adds a message layer on top of the v8worker package (binding to Javascript V8).

The https://github.com/ry/v8worker package is a simple binding that provides a few Javascript operations to send simple messages (string) to and receive from Go.
Recently, the v8worker package has been enhanced to support synchronous communication; this allows for accessing return values from functions.
The v8dispatcher package sends MessageSend values serialized as JSON strings to be dispatched in Go or Javascript.
A MessageDispatcher is used to dispatch MessageSend values to function calls, both in Go and in Javascript (no reflection).

Methods available in Go to invoke custom functions in Javascript (see MessageDispatcher):

	// Call is an asynchronous call to Javascript and does no expect a return value.
	// CallReturn is synchronous call to Javascript and expects a return value.

Functions available in Javascript to invoke custom functions in Go (see js folder):

	// V8D.call performs a MessageSend in Go and does NOT return a value.
	// V8D.callReturn performs a MessageSend in Go and returns the value from that result.
	// V8D.callThen performs a MessageSend in Go which calls the onReturn function with the result.

Variables in Javascript can be set and get using:

	// Set will add/replace the value for a global variable in Javascript.
	// Get will return the value for the global variable in Javascript.

Registration of dispatchers

Dispatching MessageSend values to functions in Go requires the registration of handlers.
The RegisterFunc can be used to map a function name (the MessageSend receiver and/or selector) to a Go function.
Alternatively, by implementing the MessageHandler interface, the mapping of selectors will have to be implemented in the Perform method.


Dispatching strategy

In Javascript, the receiver field of a MessageSend is used to find the namespace starting at the global.
An empty receiver or "this" refers to the gobal namespace. The selector is used to lookup the function in the namespace.

In Go, the receiver field of a MessageSend is used to find a handler (MessageSendHandler) in the registry of the dispatcher.
If found, the handler's Perform method is called with the MessageSend in which the selector can be inspected.
An empty receiver will cause the dispatcher to look for a registered function (MessageSendHandlerFunc) instead.

For examples see the README.md and the tests.

(c) 2016, http://ernestmicklei.com. MIT License
*/
package v8dispatcher
