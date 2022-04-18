package ssh

import (
	"cc-rshell-server/model"
	"cc-rshell-server/sockets/types"
	"cc-rshell-server/ssh/app"
	sshScreen "cc-rshell-server/ssh/screen"
	"encoding/hex"
	"github.com/gliderlabs/ssh"
	"log"
	"strconv"
	"strings"
)

func ListenAndServer(addr string, registry types.ClientRegistry) error {
	ssh.Handle(func(s ssh.Session) {
		var pubKeyStr string
		if pubKey := s.PublicKey(); pubKey != nil {
			b := strings.Builder{}
			b.WriteString(pubKey.Type())
			b.WriteString(":")
			b.WriteString(hex.EncodeToString(pubKey.Marshal()))
			pubKeyStr = b.String()
			log.Printf("[*] SSH Client (ssh_user: %s, pubkey: %s) connected from %s", s.User(), pubKeyStr, s.RemoteAddr())
		}

		log.Printf("[*] SSH Client (ssh_user: %s) connected from %s", s.User(), s.RemoteAddr())

		d := findRemoteComputer(registry, s.User())
		if d == nil {
			log.Printf("[#] Unable to find matching remote computer.")
			_, _ = s.Write([]byte("Unable to find matching remote computer.\nIs your remote machine online?\n"))
		} else {
			screen, err := sshScreen.NewSSHScreen(s)
			if err != nil {
				panic(err)
			}

			err = app.RunApp(screen, registry, s.User(), s.PublicKey())
			if err != nil {
				log.Printf("[!] Lost SSH connection to %s at %s.\n", s.User(), s.RemoteAddr())
			}
		}

		_ = s.Close()
		log.Printf("[*] SSH Client %s disconnected", s.User())
	})

	//pubKeyAuth := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
	//	return true
	//})

	log.Printf("[*] Listen for SSH connections on %s.\n", addr)
	return ssh.ListenAndServe(addr, nil, ssh.HostKeyFile("./key.pem"))
}

func findRemoteComputer(registry types.ClientRegistry, sshUser string) types.ComputerDescriptor {
	if sshUser == "" {
		return nil
	}

	if id, err := strconv.ParseInt(sshUser, 10, 64); err == nil {
		// ssh user might be a computer id
		if d, exists := registry[model.ComputerID(id)]; exists {
			return d
		}
	}

	// ssh user might be the computer label
	for _, d := range registry {
		if d.ComputerLabel() == sshUser {
			return d
		}
	}

	return nil
}
