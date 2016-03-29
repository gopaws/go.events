package events

import (
	"testing"
)

type NewListener struct {
	Target bool
}

func (handler *NewListener) Handle(event Event) {
	handler.Target = true
}

func TestString(t *testing.T) {
	name := "test"
	event := New(name)
	if event.String() != name {
		t.Log("Fail event create")
		t.Fail()
	}
}
