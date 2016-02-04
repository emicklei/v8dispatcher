/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.txt', which is part of this source code package.
 *
 * author: emicklei
 */

V8D.receiveCallback = function(msg) {
    var obj = JSON.parse(msg);
    var context = this;
    if (obj.receiver != "this") {
        var namespaces = obj.receiver.split(".");
        for (var i = 0; i < namespaces.length; i++) {
            context = context[namespaces[i]];
        }
    }
    var func = context[obj.selector];
    if (func != null) {
        return JSON.stringify(func.apply(context, obj.args));
    } else {
        // try reporting the error
        if (console != null) {
            console.log("[JS] unable to perform", msg);
        }
        // TODO return error?
        return "null";
    }
}

// This callback is set for handling function calls from Go transferred as JSON.
// It is called from Go using "worker.Send(...)".
// Throws a SyntaxError exception if the string to parse is not valid JSON.
//
$recv(V8D.receiveCallback);

// This callback is set for handling function calls from Go transferred as JSON that expect a return value.
// It is called from Go using "worker.SendSync(...)".
// Throws a SyntaxError exception if the string to parse is not valid JSON.
// Returns the JSON representation of the return value of the handling function.
//
$recvSync(V8D.receiveCallback);

// callDispatch is used from Go to call a callback function that was registered.
//
V8D.callDispatch = function(functionRef /*, arguments */ ) {
    var jsonArgs = [].slice.call(arguments).splice(1);
    var callback = V8D.function_registry.take(functionRef)
    if (V8D.function_registry.none == callback) {
        $print("[JS] no function for reference:" + functionRef);
        return;
    }	
    callback.apply(this, jsonArgs.map(function(each){ return JSON.parse(each); }));
}

// MessageSend is a constructor.
//
V8D.MessageSend = function MessageSend(receiver, selector, onReturn) {
    this.data = {
        "receiver": receiver,
        "selector": selector,
        "callback": onReturn,
        "args": [].slice.call(arguments).splice(3)
    };
}

// MessageSend toJSON returns the JSON representation.
//
V8D.MessageSend.prototype = {
    toJSON: function() {
        return JSON.stringify(this.data);
    }
}

// callReturn performs a MessageSend in Go and returns the value from that result
//
V8D.callReturn = function(receiver, selector /*, arguments */ ) {
    var msg = {
        "receiver": receiver,
        "selector": selector,
        "args": [].slice.call(arguments).splice(2)
    };
    return JSON.parse($sendSync(JSON.stringify(msg)));
}

// call performs a MessageSend in Go and does NOT return a value.
//
V8D.call = function(receiver, selector /*, arguments */ ) {
    var msg = {
        "receiver": receiver,
        "selector": selector,
        "args": [].slice.call(arguments).splice(2)
    };
    $send(JSON.stringify(msg));
}

// callThen performs a MessageSend in Go which can call the onReturn function.
// It does not return the value of the perform.
//
V8D.callThen = function(receiver, selector, onReturnFunction /*, arguments */ ) {
    var msg = {
        "receiver": receiver,
        "selector": selector,
        "callback": V8D.function_registry.put(onReturnFunction),
        "args": [].slice.call(arguments).splice(3)
    };
    $send(JSON.stringify(msg));
}

// set adds/replaces the value for a variable in the global scope.
//
V8D.set = function(variableName,itsValue) {
	V8D.outerThis[variableName] = itsValue;
}

// get returns the value for a variable in the global scope.
//
V8D.get = function(variableName) {
	return V8D.outerThis[variableName];
}