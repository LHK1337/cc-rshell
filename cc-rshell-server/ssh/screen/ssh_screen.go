package sshScreen

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
)

func NewSSHScreen(session ssh.Session) (tcell.Screen, error) {
	tty := NewSSHSessionTTY(session)

	ti, err := tcell.LookupTerminfo(tty.Term())
	if err != nil {
		return nil, err
	}

	screen, err := tcell.NewTerminfoScreenFromTtyTerminfo(tty, ti)
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	return screen, nil
}
