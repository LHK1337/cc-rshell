package sshScreen

import (
	"github.com/gliderlabs/ssh"
	"sync"
)

type SSHSessionTTY struct {
	ssh.Session

	sig   <-chan ssh.Window
	cb    func()
	stopQ chan struct{}
	dev   string
	wg    sync.WaitGroup
	l     sync.Mutex
}

func NewSSHSessionTTY(session ssh.Session) *SSHSessionTTY {
	return &SSHSessionTTY{Session: session}
}

func (tty *SSHSessionTTY) Term() string {
	pty, _, _ := tty.Pty()
	return pty.Term
}

func (tty *SSHSessionTTY) NotifyResize(cb func()) {
	tty.l.Lock()
	tty.cb = cb
	tty.l.Unlock()
}

func (tty *SSHSessionTTY) WindowSize() (width int, height int, err error) {
	pty, _, _ := tty.Pty()
	return pty.Window.Width, pty.Window.Height, nil
}

func (tty *SSHSessionTTY) Start() error {
	tty.l.Lock()
	defer tty.l.Unlock()

	tty.stopQ = make(chan struct{})
	tty.wg.Add(1)
	go func(stopQ chan struct{}) {
		defer tty.wg.Done()
		for {
			select {
			case <-tty.sig:
				tty.l.Lock()
				cb := tty.cb
				tty.l.Unlock()
				if cb != nil {
					cb()
				}
			case <-stopQ:
				return
			}
		}
	}(tty.stopQ)

	_, tty.sig, _ = tty.Pty()
	return nil
}

func (tty *SSHSessionTTY) Stop() error {
	tty.l.Lock()

	close(tty.stopQ)
	tty.l.Unlock()

	tty.wg.Wait()

	return nil
}

func (tty *SSHSessionTTY) Drain() error {
	return nil
}
