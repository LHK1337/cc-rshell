package messages

import (
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
)

func BuildCommandMessage(procID int, command string, params ...interface{}) []byte {
	bytes, _ := msgpack.Marshal(gin.H{
		"type":   "cmd",
		"procID": procID,
		"cmd":    command,
		"params": params,
	})
	return bytes
}
