package v8dispatcher

import (
	"bytes"
	"fmt"
	"log"
)

// Log can be used to inject your own logging framework
var Log = func(level, msg string, kvs ...interface{}) {
	// default uses standard logging
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "[%s] %s", level, msg)
	for i := 0; i < len(kvs); i = i + 2 {
		fmt.Fprintf(buf, ", %v = %v", kvs[i], kvs[i+1])
	}
	log.Println(buf.String)
}
