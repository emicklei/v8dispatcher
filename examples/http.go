/*

Example MessageHandler that provides an HTTP api (simple GET only) to Javascript, e.g.

	var doc = HttpAPI.get("https://api.myjson.com/bins/58is1")

To run this example:

	go run http.go

You should see:

	2016/02/08 18:00:41 [info] console.log, doc = map[title:v8dispatcher HttpAPI test]

*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	v8d "github.com/emicklei/v8dispatcher"
)

type HttpAPI struct{}

func (a HttpAPI) Perform(msg v8d.MessageSend) (interface{}, error) {
	switch msg.Selector {
	case "get":
		return a.doGet(msg)
	default:
		return nil, fmt.Errorf("unknown selector:%s", msg.Selector)
	}
}

func (a HttpAPI) doGet(msg v8d.MessageSend) (interface{}, error) {
	resp, err := http.Get(msg.Arguments[0].(string))
	if err != nil {
		return nil, err
	}
	var doc map[string]interface{}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (a HttpAPI) javascript() string {
	return `
		var HttpAPI = {};
		HttpAPI.get = function(url) {
			return V8D.callReturn("HttpAPI","get",url);
		};
	`
}

func (a HttpAPI) Register(m *v8d.MessageDispatcher) {
	m.Register("HttpAPI", a)
	m.Worker().Load("HttpAPI.js", a.javascript())
}

func main() {
	m := v8d.NewMessageDispatcher()
	a := HttpAPI{}
	a.Register(m)
	m.Worker().Load("myjson-sample.js", `
		var doc = HttpAPI.get("https://api.myjson.com/bins/58is1")
		console.log("doc",doc);	
	`)
}
