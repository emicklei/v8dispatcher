# v8dispatcher

message dispatching framework on top of v8worker  
- synchronous message calls
- console logging
- javascript message dispatch
- Go message dispatch
- function scheduler (call later)


# synchronous call from Javascript to Go
	var now = $sendSync(new V8D.MessageSend("time","Now"));

# synchronous call from Go to Javascript
	worker.SendSync(v8dispatcher.NewMessage("Date","now"));

# asynchronous call from Javascript to Go
	$send(new V8D.MessageSend("console","log","hello world"));

# asynchronous call from Go to Javascript
	worker.Send(v8dispatcher.NewMessage("receiver","",...));


# Example: console
In Javascript, you want to have

	console.log("the answer is", 42);
	
that will perform the log function of a Go counterpart

	type Console struct{}
	
	func (c Console) log(args ...interface{}) { }
	
Using this package, your Javascript source will be

	console = {};
	console.log = function() {
		$send(JSON.stringify({
			"receiver":"console",
			"selector":"log",
			"args": arguments
		}));
	};