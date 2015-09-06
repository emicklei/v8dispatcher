package v8dispatcher

import "log"

type Console struct {
}

func (c Console) Perform(msg MessageSend) (interface{}, error) {
	log.Println("msg", "call", msg.String())
	return nil, nil
}
