package v8dispatcher

import (
	"encoding/json"
	"fmt"
)

var NoMessageSend = MessageSend{}

type MessageSend struct {
	Receiver  string        `json:"receiver" `
	Method    string        `json:"method" `
	Arguments []interface{} `json:"args" `
	Callback  string        `json:"callback" `
	OnError   string        `json:"onError" `
	Stack     string        `json:"stack" `
}

func (m MessageSend) String() string {
	return fmt.Sprintf("%s.%s(%v) => (%s, %s)", m.Receiver, m.Method, m.Arguments, m.Callback, m.OnError)
}

func (m MessageSend) JSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
