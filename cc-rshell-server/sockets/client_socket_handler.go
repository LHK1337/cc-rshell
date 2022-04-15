package sockets

import (
	"cc-rshell-server/sockets/messages"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

func NewClientSocketHandler() *melody.Melody {
	s := melody.New()
	// TODO: better origin check
	s.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	s.HandleMessage(messages.MessageTransformer)

	return s
}
