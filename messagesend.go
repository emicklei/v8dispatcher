package v8dispatcher

import "encoding/json"

type MessageSend struct {
	Receiver  string        `json:"receiver" `
	Method    string        `json:"method" `
	Arguments []interface{} `json:"args" `
}

func (m MessageSend) JSON() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
