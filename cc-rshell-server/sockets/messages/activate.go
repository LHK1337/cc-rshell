package messages

import (
	"cc-rshell-server/sockets/types"
	"github.com/gin-gonic/gin"
)

type activateMessage struct {
	Id       types.ComputerID  `json:"id"`
	Label    string            `json:"label"`
	KeyCodes types.KeyCodesMap `json:"keyCodes"`
}

func handleActivateMessage(d types.ComputerDescriptor, msg gin.H) error {
	var activateMsg activateMessage
	err := parseJSONStruct(msg, &activateMsg)
	if err != nil {
		return err
	}

	d.Activate(activateMsg.Id, activateMsg.Label, activateMsg.KeyCodes)
	return nil
}
