package v8dispatcher

import (
	"bytes"
	"fmt"
	"log"
)

var Debug = false

// Log can be used to inject your own logging framework
var Log = func(level, msg string, kvs ...interface{}) {
	// default uses standard logging
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "[%s] %s", level, msg)
	for i := 0; i < len(kvs); i = i + 2 {
		var v interface{}
		if len(kvs) == i+1 {
			v = "*** missing ***"
		} else {
			v = kvs[i+1]
		}
		fmt.Fprintf(buf, ", %v = %v", kvs[i], v)
	}
	log.Println(buf.String())
}

func ConsoleLog(msg MessageSend) (interface{}, error) {
	log.Println(msg.Arguments...)
	return nil, nil
}
