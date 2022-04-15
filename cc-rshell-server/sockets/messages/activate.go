package messages

import (
	"cc-rshell-server/sockets/types"
	"github.com/gin-gonic/gin"
	"log"
)

type activateMessage struct {
	Id       types.ComputerID  `json:"id" msgpack:"id"`
	Label    string            `json:"label" msgpack:"label"`
	KeyCodes types.KeyCodesMap `json:"keyCodes" msgpack:"keyCodes"`
}

func handleActivateMessage(d types.ComputerDescriptor, msg gin.H) error {
	var activateMsg activateMessage
	err := parseDynamicStruct(msg, &activateMsg)
	if err != nil {
		return err
	}

	d.Activate(activateMsg.Id, activateMsg.Label, activateMsg.KeyCodes)

	log.Printf("[*] Client (%d:'%s') at %s activated.\n", activateMsg.Id, activateMsg.Label, d.RemoteAddr())

	return nil
}
