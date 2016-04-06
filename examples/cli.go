/*

Command line program that interacts with a V8 Javascript engine through a v8dispatcher.MessageDispather

	go run cli.go

Currently, every expression entered is wrapped in a console.log(...) call to print the value of that expression.

Append the "-v" option to see the MessageSends exchanged.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/GeertJohan/go.linenoise"
	"github.com/emicklei/v8dispatcher"
)

const historyFile = ".v8dispatcher-history"

var (
	lastHistoryEntry string
	v8d              *v8dispatcher.MessageDispatcher
	verbose          = flag.Bool("v", false, "verbose output")
)

func main() {
	flag.Parse()
	fmt.Println("V8D is ready")
	v8d = v8dispatcher.NewMessageDispatcher()

	// override default log
	v8d.RegisterFunc("console.log",
		func(msg v8dispatcher.MessageSend) (interface{}, error) {
			fmt.Println(msg.Arguments...)
			return nil, nil
		})

	// set debug level
	v8d.Trace(*verbose)
	loop()
}

func processLine(line string) string {
	if strings.HasPrefix(line, "console.log") {
		err := v8d.Worker().Load("line0.js", line)
		if err != nil {
			return err.Error()
		}
		return ""
	}
	// wrap expression in a console call to see the result value
	err := v8d.Worker().Load("line0.js", fmt.Sprintf("console.log(%s);", strings.TrimRight(line, ";")))
	if err != nil {
		return err.Error()
	}
	return ""
}

func loop() {
	linenoise.LoadHistory(historyFile)
	for {
		entered, err := linenoise.Line("> ")
		if err != nil {
			if err == linenoise.KillSignalError {
				os.Exit(0)
			}
			fmt.Println("unexpected error: %s", err)
			os.Exit(0)
		}
		entry := strings.TrimLeft(entered, "\t ") // without tabs,spaces
		var output string
		if entry != lastHistoryEntry {
			err = linenoise.AddHistory(entry)
			if err != nil {
				fmt.Printf("error: %s\n", entry)
			}
			lastHistoryEntry = entry
			linenoise.SaveHistory(historyFile)
		}
		output = processLine(entry)
		if len(output) > 0 {
			fmt.Println(output)
		}
	}
}
