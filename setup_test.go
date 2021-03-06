package v8dispatcher

import "testing"

type recorder struct {
	moduleName string
	source     string
	msg        *MessageSend
}

func (r *recorder) Perform(msg MessageSend) (interface{}, error) {
	r.msg = &msg
	return nil, nil
}

func expectConsoleLogArgument(t *testing.T, rec *recorder, arg interface{}) {
	if rec.msg == nil {
		t.Fatal("message not recorded")
	}
	if got, want := rec.msg.Selector, "log"; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := len(rec.msg.Arguments), 1; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
	if got, want := rec.msg.Arguments[0], arg; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}
