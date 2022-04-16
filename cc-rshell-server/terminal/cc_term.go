package terminal

import "cc-rshell-server/model"

// CCTerm describes the Computer Craft term.Redirect interface
type CCTerm interface {
	Write(text string)
	Scroll(y int)
	GetCursor() (x, y int)
	SetCursor(x, y int)
	GetCursorBlink() bool
	SetCursorBlink(doBlink bool)
	GetSize() (width, height int)
	Clear()
	ClearLine()
	GetTextColor() model.ColorID
	SetTextColor(color model.ColorID)
	GetBackgroundColor() model.ColorID
	SetBackgroundColor(color model.ColorID)
	IsColor()
	Blit(text, textColor, backgroundColor string)
	SetPaletteColor([]struct {
		id   model.ColorID
		code model.ColorCode
	})
	GetPaletteColor(id model.ColorID) model.ColorCode
}
