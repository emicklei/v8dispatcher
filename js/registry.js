/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.txt', which is part of this source code package.
 *
 * author: emicklei
 */

var V8D = V8D || {"globals":{}};

// http://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript
V8D.uuid = function() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0,
            v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// function_registry keeps identifyable (by generated id) functions
//
V8D.function_registry = {};
V8D.function_registry.none = undefined;
V8D.function_registry.put = function(func) {
    var ref = V8D.uuid();
    this[ref] = func;
    return ref;
}

// take returns the function by its reference and removes it from the registry.
//
V8D.function_registry.take = function(ref) {
    var func = this[ref];
    this[ref] = undefined;
    return func;
}
