/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.txt', which is part of this source code package.
 *
 * author: emicklei
 */

// This callback is set for handling function calls in JSON from Go.
// It is called from Go using "worker.Send(...)".
// Throws a SyntaxError exception if the string to parse is not valid JSON.
//
$recv(function(msg) {
    var obj = JSON.parse(msg);
    this[obj.selector].apply(this, obj.args)
});

// javascript_dispatch is used to directly call a Javascript function by its name.
// 
function javascript_dispatch(functionName, context /*, args */ ) {
    var args = [].slice.call(arguments).splice(2);
    var namespaces = functionName.split(".");
    var func = namespaces.pop();
    for (var i = 0; i < namespaces.length; i++) {
        context = context[namespaces[i]];
    }
    return context[func].apply(this, args);
}

// callback_dispatch is used from Go to call a callback function that was registered.
//
function callback_dispatch(functionRef /*, args */ ) {
	var args = [].slice.call(arguments).splice(1);
	var callback = V8D.function_registry.take(functionRef)
	if (V8D.function_registry.none == callback) {
		$print("no function for reference:"+functionRef);
		return;
	}
	callback.apply(this,args)
}


// go_dispatch is used in Javascript to call a Go function.
// if the Go function returns a non-nil value then the onReturn is called
//
function go_dispatch(onReturn, receiver, methodName /* args */ ) {
//	var iferror = function(reason) {	
//		var lines = stk.split("\n");
//		$print("js: "+reason);					
//		$print(lines);
//	}
    var obj = {		
		"receiver": receiver,
		"method": methodName,
		"callback": onReturn,
		"stack": new Error().stack
    };
    obj["args"] = [].slice.call(arguments).splice(3);
    $send(JSON.stringify(obj));
}




