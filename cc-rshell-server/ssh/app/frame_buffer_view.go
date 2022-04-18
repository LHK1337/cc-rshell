package app

import (
	"cc-rshell-server/model"
	"cc-rshell-server/utils"
	"context"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

type FrameBufferView struct {
	*tview.Box

	fbChannel chan *model.FrameBuffer
	ctx       context.Context

	bufferLock  sync.Mutex
	localBuffer *model.FrameBuffer
	blinkState  bool
	getColors   func() model.ColorPalette

	firstDataCallback func()
	requestDraw       func()
	channelClosed     func()
}

func NewFramebufferView(ctx context.Context, fbChannel chan *model.FrameBuffer, getColors func() model.ColorPalette, firstDataCallback, requestDraw, channelClosed func()) *FrameBufferView {
	return &FrameBufferView{
		Box:               tview.NewBox(),
		ctx:               ctx,
		fbChannel:         fbChannel,
		bufferLock:        sync.Mutex{},
		firstDataCallback: firstDataCallback,
		requestDraw:       requestDraw,
		getColors:         getColors,
		channelClosed:     channelClosed,
	}
}

func (v *FrameBufferView) Draw(screen tcell.Screen) {
	v.DrawForSubclass(screen, v.Box)

	v.bufferLock.Lock()
	defer v.bufferLock.Unlock()

	if v.localBuffer == nil {
		return
	}

	b := v.localBuffer

	localX, localY, w, h := v.Box.GetInnerRect()
	w = utils.Min(w, b.SizeX)
	h = utils.Min(h, b.SizeY)

	colors := v.getColors()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			bgBlit := model.Blit(string(b.BackgroundColor[y][x]))
			fgBlit := model.Blit(string(b.TextColor[y][x]))

			style := tcell.Style{}.
				Background(tcell.NewHexColor(int32(colors[bgBlit].ColorCode))).
				Foreground(tcell.NewHexColor(int32(colors[fgBlit].ColorCode)))

			if b.CursorX == x && b.CursorY == y && b.CursorBlink {
				style.Blink(v.blinkState)
				v.blinkState = !v.blinkState
			}

			screen.SetContent(localX+b.XOffset+x, localY+b.YOffset+y, rune(b.Text[y][x]), nil, style)
		}
	}
}

func (v *FrameBufferView) Worker() {
	for {
		select {
		case <-v.ctx.Done():
			return
		case b, isOpen := <-v.fbChannel:
			if !isOpen && v.channelClosed != nil {
				v.channelClosed()
			}

			v.bufferLock.Lock()
			v.localBuffer = b
			v.bufferLock.Unlock()

			if v.firstDataCallback != nil {
				v.firstDataCallback()
				v.firstDataCallback = nil
			}

			v.requestDraw()
		}
	}
}
