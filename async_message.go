package v8dispatcher

import "fmt"

var NoMessageSend = AsyncMessage{}

type AsyncMessage struct {
	MessageSend
	Callback string `json:"callback" `
	OnError  string `json:"onError" `
	Stack    string `json:"stack" `
}

func (m AsyncMessage) String() string {
	return fmt.Sprintf("%s.%s(%v) => (%s, %s)", m.Receiver, m.Method, m.Arguments, m.Callback, m.OnError)
}
