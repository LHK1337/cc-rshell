package messages

import (
	"cc-rshell-server/model"
	"cc-rshell-server/sockets/types"
	"github.com/vmihailenco/msgpack/v5"
)

type framebufferUpdatePayload struct {
	ProcID int `json:"procID" msgpack:"procID"`

	Buffer model.FrameBuffer `json:"buffer" msgpack:"buffer"`
}

func handleFrameBufferUpdateMessage(d types.ComputerDescriptor, msg []byte) error {
	channelMap := d.FramebufferChannelMap()

	var fb framebufferUpdatePayload
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

type framebufferClosedPayload struct {
	ProcID int `json:"procID" msgpack:"procID"`
}

func handleFrameBufferClosedMessage(d types.ComputerDescriptor, msg []byte) error {
	channelMap := d.FramebufferChannelMap()

	var fb framebufferClosedPayload
	err := msgpack.Unmarshal(msg, &fb)
	if err != nil {
		return err
	}

	if channel, exists := channelMap[fb.ProcID]; exists && channel != nil {
		func() {
			defer func() {
				// we do not care here
				_ = recover()
			}()

			delete(channelMap, fb.ProcID)
			close(channel)
		}()
	}

	return nil
}
