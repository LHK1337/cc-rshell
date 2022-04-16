package types

import (
	"cc-rshell-server/model"
	"gopkg.in/olahol/melody.v1"
	"time"
)

type ComputerDescriptor interface {
	// Init initializes new connections
	Init()
	// Activate activates a connection with given information
	Activate(id model.ComputerID, label string, keyCodes model.KeyCodesMap, colors model.ColorPalette)
	// Activated returns true whether the connection is activated otherwise false
	Activated() bool
	// ComputerID returns the ComputerCraft computer ID of the remote computer
	// In the scope of a minecraft world is this ID unique
	ComputerID() model.ComputerID
	// ComputerLabel returns the ComputerCraft computer label of the remote computer
	// Might NOT be unique
	ComputerLabel() string
	// RemoteAddr returns the connection's remote address
	RemoteAddr() string
	// Colors returns the current color palette
	Colors() model.ColorPalette
	// KeyCodes returns the ComputerCraft keys table with key codes used in key events by the computer
	KeyCodes() model.KeyCodesMap
	// ConnectedSince time when the connection was established
	ConnectedSince() time.Time
	// MessageBufferMap returns the buffer map of this connection
	MessageBufferMap() model.BufferMap
	// Close closes the connection
	Close() error
}

func WrapSession(s *melody.Session) ComputerDescriptor {
	return &ComputerDescriptorImpl{s}
}

const (
	SessionActivationTimeout = 10 * time.Second

	InvalidComputerID model.ComputerID = -1

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
	setValue(d, computerBufferMapKey, model.BufferMap{})

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

func (d *ComputerDescriptorImpl) Activate(id model.ComputerID, label string, keyCodes model.KeyCodesMap, colors model.ColorPalette) {
	setValue(d, computerIDKey, id)
	setValue(d, computerLabelKey, label)
	setValue(d, computerKeyCodesKey, keyCodes)
	setValue(d, computerColorsKey, colors)
	setValue(d, computerActivatedKey, true)
}

func (d *ComputerDescriptorImpl) Activated() bool {
	return getValue(d, computerActivatedKey, false)
}

func (d *ComputerDescriptorImpl) ComputerID() model.ComputerID {
	return getValue(d, computerIDKey, InvalidComputerID)
}

func (d *ComputerDescriptorImpl) RemoteAddr() string {
	return d.Request.RemoteAddr
}

func (d *ComputerDescriptorImpl) KeyCodes() model.KeyCodesMap {
	return getValue(d, computerKeyCodesKey, model.KeyCodesMap(nil))
}

func (d *ComputerDescriptorImpl) Colors() model.ColorPalette {
	return getValue(d, computerColorsKey, model.ColorPalette(nil))
}

func (d *ComputerDescriptorImpl) ComputerLabel() string {
	return getValue(d, computerLabelKey, "")
}

func (d *ComputerDescriptorImpl) MessageBufferMap() model.BufferMap {
	return getValue(d, computerBufferMapKey, model.BufferMap(nil))
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
