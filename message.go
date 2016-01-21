package v8dispatcher

import (
	"encoding/json"
	"fmt"
)

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
	return fmt.Sprintf("%s(%v) => (%s, %s)", m.Selector, m.Arguments, m.Callback, m.OnError)
}
