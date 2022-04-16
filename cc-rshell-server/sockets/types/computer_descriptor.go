package types

import (
	"gopkg.in/olahol/melody.v1"
	"time"
)

type ComputerID int
type KeyCodesMap map[string]interface{}
type BufferMap map[byte]*TimedBuffer

type Blit string
type CCColor struct {
	Label     string `json:"label" msgpack:"label"`
	ColorID   uint   `json:"colorID" msgpack:"colorID"`
	ColorCode uint   `json:"colorCode" msgpack:"colorCode"`
}

type ColorPalette map[Blit]CCColor

type ComputerDescriptor interface {
	// Init initializes new connections
	Init()
	// Activate activates a connection with given information
	Activate(id ComputerID, label string, keyCodes KeyCodesMap, colors ColorPalette)
	// Activated returns true whether the connection is activated otherwise false
	Activated() bool
	// ComputerID returns the ComputerCraft computer ID of the remote computer
	// In the scope of a minecraft world is this ID unique
	ComputerID() ComputerID
	// ComputerLabel returns the ComputerCraft computer label of the remote computer
	// Might NOT be unique
	ComputerLabel() string
	// RemoteAddr returns the connection's remote address
	RemoteAddr() string
	// Colors returns the current color palette
	Colors() ColorPalette
	// KeyCodes returns the ComputerCraft keys table with key codes used in key events by the computer
	KeyCodes() KeyCodesMap
	// ConnectedSince time when the connection was established
	ConnectedSince() time.Time
	// MessageBufferMap returns the buffer map of this connection
	MessageBufferMap() BufferMap
	// Close closes the connection
	Close() error
}

func WrapSession(s *melody.Session) ComputerDescriptor {
	return &ComputerDescriptorImpl{s}
}

const (
	SessionActivationTimeout = 10 * time.Second

	InvalidComputerID ComputerID = -1

	computerIDKey        = "CLIENT_COMPUTER_ID"
	computerLabelKey     = "CLIENT_COMPUTER_LABEL"
	computerKeyCodesKey  = "CLIENT_COMPUTER_KEY_CODES"
	computerColorsKey    = "CLIENT_COMPUTER_COLORS"
	computerBufferMapKey = "CLIENT_COMPUTER_BUFFER_MAP"
	computerActivatedKey = "CLIENT_COMPUTER_ACTIVATED"
	connectedSinceKey    = "CLIENT_COMPUTER_CONNECTED_SINCE"
)

type ComputerDescriptorImpl struct {
	*melody.Session
	// DO NOT PUT PROPERTIES HERE
	// Store them in the Session with getValue() and setValue()
}

func (d *ComputerDescriptorImpl) Init() {
	d.initTimeout(SessionActivationTimeout, d.Close)
}

// Helps in tests
func (d *ComputerDescriptorImpl) initTimeout(timout time.Duration, closeFunc func() error) {
	setValue(d, connectedSinceKey, time.Now())
	setValue(d, computerActivatedKey, false)
	setValue(d, computerBufferMapKey, BufferMap{})

	d.startActivationTimeout(timout, closeFunc)
}

// Helps in tests
func (d *ComputerDescriptorImpl) startActivationTimeout(timout time.Duration, closeFunc func() error) {
	go func() {
		// Silence fail
		defer func() {
			_ = recover()
		}()

		<-time.After(timout)
		if !d.Activated() {
			_ = closeFunc()
		}
	}()
}

func (d *ComputerDescriptorImpl) Activate(id ComputerID, label string, keyCodes KeyCodesMap, colors ColorPalette) {
	setValue(d, computerIDKey, id)
	setValue(d, computerLabelKey, label)
	setValue(d, computerKeyCodesKey, keyCodes)
	setValue(d, computerColorsKey, colors)
	setValue(d, computerActivatedKey, true)
}

func (d *ComputerDescriptorImpl) Activated() bool {
	return getValue(d, computerActivatedKey, false)
}

func (d *ComputerDescriptorImpl) ComputerID() ComputerID {
	return getValue(d, computerIDKey, InvalidComputerID)
}

func (d *ComputerDescriptorImpl) RemoteAddr() string {
	return d.Request.RemoteAddr
}

func (d *ComputerDescriptorImpl) KeyCodes() KeyCodesMap {
	return getValue(d, computerKeyCodesKey, KeyCodesMap(nil))
}

func (d *ComputerDescriptorImpl) Colors() ColorPalette {
	return getValue(d, computerColorsKey, ColorPalette(nil))
}

func (d *ComputerDescriptorImpl) ComputerLabel() string {
	return getValue(d, computerLabelKey, "")
}

func (d *ComputerDescriptorImpl) MessageBufferMap() BufferMap {
	return getValue(d, computerBufferMapKey, BufferMap(nil))
}

func (d *ComputerDescriptorImpl) ConnectedSince() time.Time {
	return getValue(d, connectedSinceKey, time.Time{})
}

func (d *ComputerDescriptorImpl) Close() error {
	setValue(d, computerActivatedKey, false)

	// do not care if it is already closed
	_ = d.Session.Close()

	return nil
}

func getValue[T any](d *ComputerDescriptorImpl, key string, elseValue T) T {
	valueInterface, exists := d.Keys[key]
	if value, ok := valueInterface.(T); ok && exists {
		return value
	}
	return elseValue
}

func setValue[T any](d *ComputerDescriptorImpl, key string, value T) {
	d.Keys[key] = value
}
