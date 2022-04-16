package messages

import (
	"cc-rshell-server/sockets/types"
	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/olahol/melody.v1"
	"log"
	"time"
)

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

	MessageHandler(types.WrapSession(session), bytes, r)
}

type baseMessage struct {
	Type string `json:"type" msgpack:"type"`
}

func MessageHandler(d types.ComputerDescriptor, msg []byte, r types.ClientRegistry) {
	var baseMSG baseMessage
	err := msgpack.Unmarshal(msg, &baseMSG)
	if err != nil {
		log.Printf("[!] Received invalid MessagePack from client at %s.\n", d.RemoteAddr())
		return
	}

	if baseMSG.Type == "" {
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
			baseMSG.Type, d.ComputerID(), d.ComputerLabel(), d.RemoteAddr())
	} else {
		log.Printf("[*] Received %s message from client at %s.\n", baseMSG.Type, d.RemoteAddr())
	}

	switch baseMSG.Type {
	case "activate":
		err = handleActivateMessage(d, msg)
		if err == nil {
			r[d.ComputerID()] = d
		}
	}

	if err != nil {
		if d.Activated() {
			log.Printf("[!] Unable to handle %s message from client (%d:'%s') at %s. Error: %s\n",
				baseMSG.Type, d.ComputerID(), d.ComputerLabel(), d.RemoteAddr(), err)
		} else {
			log.Printf("[!] Unable to handle %s message from client at %s. Error: %s\n",
				baseMSG.Type, d.RemoteAddr(), err)
		}
		return
	}
}
