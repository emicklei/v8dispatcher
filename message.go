package v8dispatcher

import (
	"encoding/json"
	"fmt"
)

// MessageSend encapsulates a performable message between Javascript and Go
type MessageSend struct {
	// Receiver is used to lookup a registered MesageHandler
	Receiver string `json:"receiver" `

	// Selector is used to lookup a function or statement in a registered MessageHandler.
	Selector string `json:"selector" `

	// Arguments holds any arguments needed for the function
	Arguments []interface{} `json:"args" `

	// Callback can be a function reference registered in Javascript to call back into Javascript.
	Callback string `json:"callback" `

	// IsAsynchronous is to used to indicate that no return value is expected
	IsAsynchronous bool `json:"async"`
}

func (m MessageSend) JSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// String exists for debugging
func (m MessageSend) String() string {
	return fmt.Sprintf("%#v", m)
}
