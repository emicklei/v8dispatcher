package v8dispatcher

import "log"

type Console struct {
}

func (c Console) Perform(msg AsyncMessage) (interface{}, error) {
	log.Println("msg", "call", msg.String())
	return nil, nil
}
