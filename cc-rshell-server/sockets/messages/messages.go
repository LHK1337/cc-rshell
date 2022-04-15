package messages

import (
	"cc-rshell-server/sockets/types"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/olahol/melody.v1"
	"log"
	"reflect"
)

func parseDynamicStruct[T any](dynStruct gin.H, value *T) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Unable to parse dynamic struct. Reason: %v", r))
		}
	}()

	v := reflect.Indirect(reflect.ValueOf(value))
	valueType := v.Type()
	fields := valueType.NumField()
	for i := 0; i < fields; i++ {
		field := valueType.Field(i)

		dynKey := field.Tag.Get("msgpack")
		if dynKey == "" || dynKey == "-" {
			continue
		}

		fieldValue, exists := dynStruct[dynKey]
		if !exists {
			return errors.New("Missing required field " + dynKey)
		}

		fieldType := field.Type
		if reflect.TypeOf(fieldValue).ConvertibleTo(fieldType) || reflect.TypeOf(fieldValue).AssignableTo(fieldType) {
			v.Field(i).Set(reflect.ValueOf(fieldValue).Convert(fieldType))
		} else {
			return errors.New(fmt.Sprintf("Unable to assign %v to %s (%s)", fieldValue, field.Name, fieldType))
		}
	}

	return nil
}

const (
	// Messages can be chunked
	// In that case the first byte indicate the buffer state and is followed by another byte containing the buffer id.
	// Messages that are not chunked just start with a byte indicating that.
	messageNotChunkedByte          int8 = 0x00
	chunkedMessageStartByte        int8 = 0x01
	chunkedMessageIntermediateByte int8 = 0x02
	chunkedMessageEndByte          int8 = 0x03
)

func MessageTransformer(session *melody.Session, bytes []byte) {
	var msg gin.H
	err := msgpack.Unmarshal(bytes, &msg)
	if err != nil {
		log.Printf("[!] Received invalid MessagePack from client at %s.\n", session.Request.RemoteAddr)
		return
	}

	MessageHandler(types.WrapSession(session), msg)
}

func MessageHandler(d types.ComputerDescriptor, msg gin.H) {
	msgType, exists := msg["type"]
	if !exists {
		if d.Activated() {
			log.Printf("[!] Received untyped message from client (%d:'%s') at %s.\n",
				d.ComputerID(), d.ComputerLabel(), d.RemoteAddr())
		} else {
			log.Printf("[!] Received untyped message from client at %s.\n", d.RemoteAddr())
		}
		return
	}

	if d.Activated() {
		log.Printf("[*] Received %s message from client (%d:'%s') at %s.\n",
			msgType, d.ComputerID(), d.ComputerLabel(), d.RemoteAddr())
	} else {
		log.Printf("[*] Received %s message from client at %s.\n", msgType, d.RemoteAddr())
	}

	var err error
	switch msgType {
	case "activate":
		err = handleActivateMessage(d, msg)
	}

	if err != nil {
		if d.Activated() {
			log.Printf("[!] Unable to handle %s message from client (%d:'%s') at %s. Error: %s\n",
				msgType, d.ComputerID(), d.ComputerLabel(), d.RemoteAddr(), err)
		} else {
			log.Printf("[!] Unable to handle %s message from client at %s. Error: %s\n", msgType, d.RemoteAddr(), err)
		}
		return
	}
}
