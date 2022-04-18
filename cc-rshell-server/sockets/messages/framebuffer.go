package messages

import (
	"cc-rshell-server/model"
	"cc-rshell-server/sockets/types"
	"github.com/vmihailenco/msgpack/v5"
)

type framebufferPayload struct {
	ProcID int `json:"procID" msgpack:"procID"`

	Buffer model.FrameBuffer `json:"buffer" msgpack:"buffer"`
}

func handleFrameBufferMessage(d types.ComputerDescriptor, msg []byte) error {
	channelMap := d.FramebufferChannelMap()

	var fb framebufferPayload
	err := msgpack.Unmarshal(msg, &fb)
	if err != nil {
		return err
	}

	if channel, exists := channelMap[fb.ProcID]; exists && channel != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// there was an error writing to that channel -> probably closed
					// therefore we can remove this map entry
					delete(channelMap, fb.ProcID)
				}
			}()

			channel <- &fb.Buffer
		}()
	}

	return nil
}
