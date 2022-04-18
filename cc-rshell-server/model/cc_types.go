package model

type ComputerID int
type KeyCodesMap map[string]interface{}
type BufferMap map[byte]*TimedBuffer

type Blit string
type ColorID uint
type ColorCode uint
type CCColor struct {
	Label     string    `json:"label" msgpack:"label"`
	ColorID   ColorID   `json:"colorID" msgpack:"colorID"`
	ColorCode ColorCode `json:"colorCode" msgpack:"colorCode"`
}

type ColorPalette map[Blit]CCColor

type FrameBuffer struct {
	IsColor                bool     `json:"color" msgpack:"color"`
	CurrentBackgroundColor string   `json:"curBackColor" msgpack:"curBackColor"`
	CurrentTextColor       string   `json:"curTextColor" msgpack:"curTextColor"`
	CursorBlink            bool     `json:"cursorBlink" msgpack:"cursorBlink"`
	CursorX                int      `json:"cursorX" msgpack:"cursorX"`
	CursorY                int      `json:"cursorY" msgpack:"cursorY"`
	MaxX                   int      `json:"maxX" msgpack:"maxX"`
	MaxY                   int      `json:"maxY" msgpack:"maxY"`
	MinX                   int      `json:"minX" msgpack:"minX"`
	MinY                   int      `json:"minY" msgpack:"minY"`
	SizeX                  int      `json:"sizeX" msgpack:"sizeX"`
	SizeY                  int      `json:"sizeY" msgpack:"sizeY"`
	Text                   []string `json:"text" msgpack:"text"`
	TextColor              []string `json:"textColor" msgpack:"textColor"`
	BackgroundColor        []string `json:"backColor" msgpack:"backColor"`
	XOffset                int      `json:"xOffset" msgpack:"xOffset"`
	YOffset                int      `json:"yOffset" msgpack:"yOffset"`
}
