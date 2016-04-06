# v8dispatcher

v8dispatcher is a Go package for communicating to and from Javascript running in V8.

It provides a message abstraction layer on top of the v8worker package which has been enhanced to support synchronous messaging.
The v8dispatcher has a MessageDispatcher component that is used to dispatch MessageSend values to function calls, both in Go and in Javascript (no reflection).

[![GoDoc](https://godoc.org/github.com/emicklei/v8dispatcher?status.svg)](https://godoc.org/github.com/emicklei/v8dispatcher)

### Calling Go from Javascript

__Go__

	md := NewMessageDispatcher()
	md.RegisterFunc("now",func(m MessageSend) (interface{},error) {
		return time.Now(), nil	
	})
	
__Javascript__

	var now = V8D.callReturn("","now");		
	
This is the minimal example for which a simple function (no arguments) is registered and called from Javascript when loading the source like this:

__Go__	
	
	md.Worker().Load("example.js", `var now = V8D.callReturn("","now");`)


### Calling Javascript from Go

__Javascript__

	function now() {
		return new Date();
	}
	
__Go__

	md := NewMessageDispatcher()
	now, _ := md.CallReturn("this","now")
	
	
### Asynchronous call from Javascript

__Go__

	md := NewMessageDispatcher()
	md.RegisterFunc("handleEvent",func(m MessageSend) (interface{},error) {
		dataMap := m.Arguments[0].(map[string]interface{})
		data := dataMap["data"]
		...
		return nil, nil	
	})

__Javascript__

	V8D.call("","handleEvent", {"data": "some event data"});
	
### Asynchronous call from Go

__Javascript__

	function handleEvent(data) {
		...
	}

__Go__

	md := NewMessageDispatcher()
	md.Call("this","handleEvent",map[string]interface{}{
		"data" : "some event data",
	})
	
### Set and Get global variables

__Go__
		
	md := NewMessageDispatcher()
	md.Set("shoeSize",42)
	shoeSize, _ := md.Get("shoeSize")
	
__Javascript__

	var shoeSize = this["shoeSize"]	
	
### MessageHandler

To invoke Go methods from Javascript, you can register a value whoes type implements the `MessageHandler` interface.

__Go__

	type MusicPlayer struct {}
	
	func (m MusicPlayer) Perform(m MessageSend) (interface{}, error) {
		switch (m.Selector) {
			case "start":
			case "stop":
			case "pause": 
			case "reset": 
			default: return nil , errors.New("unknown selector")
		}
	}

Register an instance of MusicPlayer

__Go__

	player := MusicPlayer{}
	md := NewMessageDispatcher()
	md.Register("player", player)

Now you can use this from Javascript

__Javascript__

	V8D.call("player","start");
	
	
(c) 2016, http://ernestmicklei.com. MIT License	