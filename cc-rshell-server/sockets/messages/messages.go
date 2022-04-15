package messages

import (
	"cc-rshell-server/sockets/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"log"
	"reflect"
)

func parseJSONStruct[T any](jsonStruct gin.H, value *T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Unable to parse json struct. Reason: %v", r))
		}
	}()

	v := reflect.Indirect(reflect.ValueOf(value))
	valueType := v.Type()
	fields := valueType.NumField()
	for i := 0; i < fields; i++ {
		field := valueType.Field(i)

		jsonKey := field.Tag.Get("json")
		if jsonKey == "" || jsonKey == "-" {
			continue
		}

		fieldValue, exists := jsonStruct[jsonKey]
		if !exists {
			return errors.New("Missing required JSON field " + jsonKey)
		}

		fieldType := field.Type
		if reflect.TypeOf(fieldValue).ConvertibleTo(fieldType) {
			v.Field(i).Set(reflect.ValueOf(fieldValue).Convert(fieldType))
		} else {
			return errors.New(fmt.Sprintf("Unable to assign %v to %s (%s)", fieldValue, field.Name, fieldType))
		}
	}

	return nil
}

func MessageTransformer(session *melody.Session, bytes []byte) {
	var msg gin.H
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		log.Println("[!] Received invalid JSON from client.")
		return
	}

	MessageHandler(types.WrapSession(session), msg)
}

func MessageHandler(d types.ComputerDescriptor, msg gin.H) {
	msgType, exists := msg["type"]
	if !exists {
		log.Println("[!] Received untyped message client.")
		return
	}

	var err error
	switch msgType {
	case "activate":
		err = handleActivateMessage(d, msg)
	}

	if err != nil {
		log.Printf("[!] Unable to handle %s message.\n", msgType)
		return
	}
}
