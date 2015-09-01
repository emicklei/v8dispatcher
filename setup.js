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
    this[obj.method].apply(this, obj.args)
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
	var callback = function_registry.take(functionRef)
	if (undefined == callback) {
		$print("no function for reference:"+functionRef);
		return;
	}
	callback.apply(this,args)
}

// go_dispatch is used in Javascript to call a Go function.
// the worker callback in Go will dispatch a MessageSend (unmarshalled from the JSON message).
//
function go_dispatch(callbackOrNull, receiver, methodName /* args */ ) {
    var obj = {
		"callback": callbackOrNull,
        "receiver": receiver,
        "method": methodName
    };
    obj["args"] = [].slice.call(arguments).splice(3);
    $send(JSON.stringify(obj));
}

// http://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript
function uuid() {
	return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
   		var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
   		return v.toString(16);
	});
}
	
// function_registry keeps identifyable (by generated id) functions
//
function_registry = {};
function_registry.put = function(func){
	var ref = uuid();
	function_registry[ref] = func;
	return ref;
}

// take returns the function by its reference and removed it from the registry.
//
function_registry.take = function(ref) {
	var func = function_registry[ref];
	function_registry[ref] = undefined;
	return func;
}

// console is used for getting log entries in a logger on the Go side.
//
console = {};
console.print = function(args) {
    var msg = "";
    for (var i = 0; i < arguments.length; i++) {
        msg += arguments[i] + " (" + typeof(arguments[i]) + ") ";
    }
    $print(msg)
}
console.log = function(args) {
    go_dispatch(undefined, "console", "log", args);
}