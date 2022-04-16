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
	"time"
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
	// Messages have a chunk header byte and can be chunked
	// The first two bits indicate the buffer op code and is followed by 6 bits containing the buffer id.
	messageNotChunkedByte          byte = 0x00
	chunkedMessageStartByte        byte = 0x01
	chunkedMessageIntermediateByte byte = 0x02
	chunkedMessageEndByte          byte = 0x03

	maxTotalBufferSizePerUser = 1 * 1024 * 1024 * 1024 // 1MB
)

func MessageTransformer(session *melody.Session, bytes []byte, r types.ClientRegistry) {
	if len(bytes) == 0 {
		return
	}

	chunkHeader := bytes[0]
	opCode := (chunkHeader & 0b11000000) >> 6
	bufID := chunkHeader & 0b00111111
	bytes = bytes[1:]
	if opCode != messageNotChunkedByte {
		if len(bytes) == 0 {
			return
		}

		d := types.WrapSession(session)
		bm := d.MessageBufferMap()

		totalBytes := 0
		for _, buffer := range bm {
			totalBytes += buffer.Buffer.Len()
			if totalBytes > maxTotalBufferSizePerUser {
				if d.Activated() {
					log.Printf("[!] Client (%d:'%s') at %s exceeded its buffer for chunked messages. "+
						"Disconnecting...\n", d.ComputerID(), d.ComputerLabel(), session.Request.RemoteAddr)
				} else {
					log.Printf("[!] Client at %s exceeded its buffer for chunked messages. Disconnecting...\n",
						session.Request.RemoteAddr)
				}

				_ = d.Close()
				return
			}
		}

		switch opCode {
		case chunkedMessageStartByte:
			b, exists := bm[bufID]
			if !exists {
				b = types.NewTimedBuffer()
				bm[bufID] = b
			}

			b.Lock.Lock()
			b.Buffer.Reset()
			b.Buffer.Write(bytes)
			b.LastModification = time.Now()
			b.Lock.Unlock()

			return
		case chunkedMessageIntermediateByte:
			b, exists := bm[bufID]
			if !exists {
				log.Printf("[!] Client at %s tried to write to an nonexistent buffer.\n", session.Request.RemoteAddr)
				return
			}

			b.Lock.Lock()
			b.Buffer.Write(bytes)
			b.LastModification = time.Now()
			b.Lock.Unlock()

			return
		case chunkedMessageEndByte:
			b, exists := bm[bufID]
			if !exists {
				log.Printf("[!] Client at %s tried to write to an nonexistent buffer.\n", session.Request.RemoteAddr)
				return
			}

			b.Lock.Lock()
			b.Buffer.Write(bytes)
			newBytes := make([]byte, b.Buffer.Len())
			copy(newBytes, b.Buffer.Bytes())
			bytes = newBytes
			b.Buffer.Reset()
			b.LastModification = time.Now()
			b.Lock.Unlock()
		default:
			log.Printf("[!] Client at %s send an invalid buffer state.\n", session.Request.RemoteAddr)
			return
		}
	}

	var msg gin.H
	err := msgpack.Unmarshal(bytes, &msg)
	if err != nil {
		log.Printf("[!] Received invalid MessagePack from client at %s.\n", session.Request.RemoteAddr)
		return
	}

	MessageHandler(types.WrapSession(session), msg, r)
}

func MessageHandler(d types.ComputerDescriptor, msg gin.H, r types.ClientRegistry) {
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
		if err == nil {
			r[d.ComputerID()] = d
		}
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
