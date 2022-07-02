package ui

import (
	"strings"
)

var Nil NilMessage

type Message interface {
	String() string
}

type InfoMsg struct {
	Message string
}

func (m InfoMsg) String() string {
	return m.Message
}

type NilMessage struct{}

func (m NilMessage) String() string {
	return strings.Repeat(".", 50)
}

type ErrorMsg struct{ Error error }

func (e ErrorMsg) String() string {
	return e.Error.Error()
}
