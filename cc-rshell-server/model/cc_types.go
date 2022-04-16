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
