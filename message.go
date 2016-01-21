package v8dispatcher

import "encoding/json"

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
