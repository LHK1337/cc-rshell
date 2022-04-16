package messages

import (
	"cc-rshell-server/sockets/types"
	"github.com/vmihailenco/msgpack/v5"
	"log"
)

type activateMessage struct {
	Id       types.ComputerID   `json:"id" msgpack:"id"`
	Label    string             `json:"label" msgpack:"label"`
	KeyCodes types.KeyCodesMap  `json:"keyCodes" msgpack:"keyCodes"`
	Colors   types.ColorPalette `json:"colors" msgpack:"colors"`
}

func handleActivateMessage(d types.ComputerDescriptor, msg []byte) error {
	var activateMsg activateMessage
	err := msgpack.Unmarshal(msg, &activateMsg)
	if err != nil {
		return err
	}

	d.Activate(activateMsg.Id, activateMsg.Label, activateMsg.KeyCodes, activateMsg.Colors)

	log.Printf("[*] Client (%d:'%s') at %s activated.\n", activateMsg.Id, activateMsg.Label, d.RemoteAddr())

	return nil
}
