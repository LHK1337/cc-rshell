package messages

import (
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
)

func BuildCommandMessage(procID, bufferWidth, bufferHeight int, command string) []byte {
	bytes, _ := msgpack.Marshal(gin.H{
		"type":   "cmd",
		"procID": procID,
		"cmd":    command,
		"bufW":   bufferWidth,
		"bufH":   bufferHeight,
	})
	return bytes
}

func BuildKillMessage(procID int) []byte {
	bytes, _ := msgpack.Marshal(gin.H{
		"type":   "kill",
		"procID": procID,
	})
	return bytes
}
