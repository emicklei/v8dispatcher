# v8dispatcher

v8worker add-on for 
- synchronous message calls
- console logging
- javascript message dispatch
- Go message dispatch
- function scheduler (call later)



# console
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