package v8dispatcher

import (
	"encoding/json"
	"fmt"
)

var NoMessageSend = MessageSend{}

type MessageSend struct {
	Receiver       string        `json:"receiver" `
	Selector       string        `json:"selector" `
	Arguments      []interface{} `json:"args" `
	IsAsynchronous bool          `json:"async"`
}

func (m MessageSend) JSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type AsyncMessage struct {
	MessageSend
	Callback string `json:"callback" `
	OnError  string `json:"onError" `
	Stack    string `json:"stack" `
}

func (m AsyncMessage) String() string {
	return fmt.Sprintf("%s.%s(%v) => (%s, %s)", m.Receiver, m.Selector, m.Arguments, m.Callback, m.OnError)
}
