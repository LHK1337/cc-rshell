package sockets

import (
	"cc-rshell-server/sockets/messages"
	"cc-rshell-server/sockets/types"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
)

type ClientSocketHandler struct {
	*melody.Melody
	types.ClientRegistry
}

func NewClientSocketHandler() *ClientSocketHandler {
	s := melody.New()
	r := types.ClientRegistry{}

	// TODO: better origin check
	s.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	s.HandleConnect(func(session *melody.Session) {
		d := types.WrapSession(session)
		d.Init()
		log.Printf("[*] CLient connected from %s\n", d.RemoteAddr())
	})
	s.HandleClose(func(session *melody.Session, _ int, _ string) error {
		return types.WrapSession(session).Close()
	})
	s.HandleDisconnect(func(session *melody.Session) {
		d := types.WrapSession(session)
		if d.Activated() {
			delete(r, d.ComputerID())
			log.Printf("[*] CLient (%d:'%s') at %s disconnected\n", d.ComputerID(), d.ComputerLabel(), d.RemoteAddr())
		} else {
			log.Printf("[*] CLient (unactivated) at %s disconnected\n", d.RemoteAddr())
		}
	})
	s.HandleMessageBinary(messages.MessageTransformer)

	return &ClientSocketHandler{s, r}
}
