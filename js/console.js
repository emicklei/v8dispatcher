/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.txt', which is part of this source code package.
 *
 * author: emicklei
 */
// console is used for getting log entries in a logger on the Go side.
//
console = {};
console.print = function() {
    var msg = "";
    for (var i = 0; i < arguments.length; i++) {
        msg += arguments[i] + " (" + typeof(arguments[i]) + ") ";
    }
    $print(msg)
}
// log takes a variable number of arguments
//
console.log = function() {
	var args = [];
	// flatten all arguments	
	for (var i = 0; i < arguments.length; i++) {
       args.push(arguments[i]);
    }
	$send(JSON.stringify({
		"receiver":"console",
		"selector":"log",
		"args": args
	}));
}