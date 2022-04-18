package messages

import (
	"cc-rshell-server/model"
	"github.com/gdamore/tcell/v2"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack/v5"
)

func BuildEventMessage(event string, params ...interface{}) []byte {
	bytes, _ := msgpack.Marshal(gin.H{
		"type":   "event",
		"event":  event,
		"params": params,
	})
	return bytes
}

func MapToCCEvents(e *tcell.EventKey, keys model.KeyCodesMap) (messages [][]byte) {
	if e.Key() == tcell.KeyRune {
		return [][]byte{BuildEventMessage("char", string(e.Rune()))}
	}

	var keyUpStack [][]byte

	setKey := func(k interface{}) {
		messages = append(messages, BuildEventMessage("key", k, false))
		keyUpStack = append(keyUpStack, BuildEventMessage("key_up", k))
	}

	if e.Modifiers()&tcell.ModCtrl > 0 {
		if k, ok := keys["leftCtrl"]; ok {
			setKey(k)
		}
	}

	if e.Modifiers()&tcell.ModAlt > 0 {
		if k, ok := keys["leftAlt"]; ok {
			setKey(k)
		}
	}

	if e.Modifiers()&tcell.ModShift > 0 {
		if k, ok := keys["leftShift"]; ok {
			setKey(k)
		}
	}

	if e.Modifiers()&tcell.ModMeta > 0 {
		if k, ok := keys["leftSuper"]; ok {
			setKey(k)
		}
	}

	keyID := lookupKey(e.Key())
	if keyID == "" {
		return nil
	}

	if k, ok := keys[keyID]; ok {
		setKey(k)
	}

	for i := len(keyUpStack) - 1; i >= 0; i-- {
		messages = append(messages, keyUpStack[i])
	}

	return messages
}

func lookupKey(key tcell.Key) string {
	switch key {
	case tcell.KeyEnter:
		return "enter"
	case tcell.KeyBackspace:
		return "backspace"
	case tcell.KeyTab:
		return "tab"
	case tcell.KeyBacktab:
		return "backTab"
	case tcell.KeyEsc:
		return "esc"
	case tcell.KeyBackspace2:
		return "backspace"
	case tcell.KeyDelete:
		return "delete"
	case tcell.KeyInsert:
		return "insert"
	case tcell.KeyUp:
		return "up"
	case tcell.KeyDown:
		return "down"
	case tcell.KeyLeft:
		return "left"
	case tcell.KeyRight:
		return "right"
	case tcell.KeyHome:
		return "home"
	case tcell.KeyEnd:
		return "end"
	//case tcell.KeyUpLeft:         return "UpLeft"
	//case tcell.KeyUpRight:        return "UpRight"
	//case tcell.KeyDownLeft:       return "DownLeft"
	//case tcell.KeyDownRight:      return "DownRight"
	//case tcell.KeyCenter:         return "Center"
	case tcell.KeyPgDn:
		return "pageDown"
	case tcell.KeyPgUp:
		return "pageUp"
	//case tcell.KeyClear:          return "Clear"
	//case tcell.KeyExit:           return "Exit"
	//case tcell.KeyCancel:         return "Cancel"
	case tcell.KeyPause:
		return "pause"
	case tcell.KeyPrint:
		return "printScreen"
	case tcell.KeyF1:
		return "f1"
	case tcell.KeyF2:
		return "f2"
	case tcell.KeyF3:
		return "f3"
	case tcell.KeyF4:
		return "f4"
	case tcell.KeyF5:
		return "f5"
	case tcell.KeyF6:
		return "f6"
	case tcell.KeyF7:
		return "f7"
	case tcell.KeyF8:
		return "f8"
	case tcell.KeyF9:
		return "f9"
	case tcell.KeyF10:
		return "f10"
	case tcell.KeyF11:
		return "f11"
	case tcell.KeyF12:
		return "f12"
	case tcell.KeyF13:
		return "f13"
	case tcell.KeyF14:
		return "f14"
	case tcell.KeyF15:
		return "f15"
	case tcell.KeyF16:
		return "f16"
	case tcell.KeyF17:
		return "f17"
	case tcell.KeyF18:
		return "f18"
	case tcell.KeyF19:
		return "f19"
	case tcell.KeyF20:
		return "f20"
	case tcell.KeyF21:
		return "f21"
	case tcell.KeyF22:
		return "f22"
	case tcell.KeyF23:
		return "f23"
	case tcell.KeyF24:
		return "f24"
	case tcell.KeyF25:
		return "f25"
	case tcell.KeyF26:
		return "f26"
	case tcell.KeyF27:
		return "f27"
	case tcell.KeyF28:
		return "f28"
	case tcell.KeyF29:
		return "f29"
	case tcell.KeyF30:
		return "f30"
	case tcell.KeyF31:
		return "f31"
	case tcell.KeyF32:
		return "f32"
	case tcell.KeyF33:
		return "f33"
	case tcell.KeyF34:
		return "f34"
	case tcell.KeyF35:
		return "f35"
	case tcell.KeyF36:
		return "f36"
	case tcell.KeyF37:
		return "f37"
	case tcell.KeyF38:
		return "f38"
	case tcell.KeyF39:
		return "f39"
	case tcell.KeyF40:
		return "f40"
	case tcell.KeyF41:
		return "f41"
	case tcell.KeyF42:
		return "f42"
	case tcell.KeyF43:
		return "f43"
	case tcell.KeyF44:
		return "f44"
	case tcell.KeyF45:
		return "f45"
	case tcell.KeyF46:
		return "f46"
	case tcell.KeyF47:
		return "f47"
	case tcell.KeyF48:
		return "f48"
	case tcell.KeyF49:
		return "f49"
	case tcell.KeyF50:
		return "f50"
	case tcell.KeyF51:
		return "f51"
	case tcell.KeyF52:
		return "f52"
	case tcell.KeyF53:
		return "f53"
	case tcell.KeyF54:
		return "f54"
	case tcell.KeyF55:
		return "f55"
	case tcell.KeyF56:
		return "f56"
	case tcell.KeyF57:
		return "f57"
	case tcell.KeyF58:
		return "f58"
	case tcell.KeyF59:
		return "f59"
	case tcell.KeyF60:
		return "f60"
	case tcell.KeyF61:
		return "f61"
	case tcell.KeyF62:
		return "f62"
	case tcell.KeyF63:
		return "f63"
	case tcell.KeyF64:
		return "f64"
	default:
		return ""
	}
}
