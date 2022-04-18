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
	// RegisterFramebufferChannel allows to register a channel to receive framebuffer changes
	RegisterFramebufferChannel(procID int, framebufferChannel chan *model.FrameBuffer)
	// FramebufferChannelMap returns a map containing framebuffer channel for all procIDs on this computer
	FramebufferChannelMap() (channelMap map[int]chan *model.FrameBuffer)
	// Close closes the connection
	Close() error
}

func WrapSession(s *melody.Session) ComputerDescriptor {
	return &ComputerDescriptorImpl{s}
}

const (
	SessionActivationTimeout = 10 * time.Second

	InvalidComputerID model.ComputerID = -1

	computerStateKey = "CLIENT_COMPUTER_STATE"
)

type ComputerDescriptorImpl struct {
	*melody.Session
	// DO NOT PUT PROPERTIES HERE
	// Store them in the computerState
}

type computerState struct {
	Activated             bool
	ID                    model.ComputerID
	Label                 string
	KeyCodes              model.KeyCodesMap
	Colors                model.ColorPalette
	FramebufferChannelMap map[int]chan *model.FrameBuffer
	MessageBufferMap      model.BufferMap
	ConnectedSince        time.Time
}

func (d *ComputerDescriptorImpl) Init() {
	d.initTimeout(SessionActivationTimeout, d.Close)
}

// Helps in tests
func (d *ComputerDescriptorImpl) initTimeout(timout time.Duration, closeFunc func() error) {
	setValue(d, computerStateKey, &computerState{
		ID:                    InvalidComputerID,
		Activated:             false,
		MessageBufferMap:      model.BufferMap{},
		ConnectedSince:        time.Now(),
		FramebufferChannelMap: map[int]chan *model.FrameBuffer{},
	})

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

func (d *ComputerDescriptorImpl) RegisterFramebufferChannel(procID int, framebufferChannel chan *model.FrameBuffer) {
	state := d.state()
	if state.FramebufferChannelMap[procID] != nil {
		close(state.FramebufferChannelMap[procID])
	}
	state.FramebufferChannelMap[procID] = framebufferChannel
}

func (d *ComputerDescriptorImpl) FramebufferChannelMap() (channelMap map[int]chan *model.FrameBuffer) {
	return d.state().FramebufferChannelMap
}

func (d *ComputerDescriptorImpl) Activate(id model.ComputerID, label string, keyCodes model.KeyCodesMap, colors model.ColorPalette) {
	state := d.state()
	state.ID = id
	state.Label = label
	state.KeyCodes = keyCodes
	state.Colors = colors
	state.Activated = true
}

func (d *ComputerDescriptorImpl) state() *computerState {
	var nullState computerState
	return getValue(d, computerStateKey, &nullState)
}

func (d *ComputerDescriptorImpl) Activated() bool {
	return d.state().Activated
}

func (d *ComputerDescriptorImpl) ComputerID() model.ComputerID {
	return d.state().ID
}

func (d *ComputerDescriptorImpl) RemoteAddr() string {
	return d.Request.RemoteAddr
}

func (d *ComputerDescriptorImpl) KeyCodes() model.KeyCodesMap {
	return d.state().KeyCodes
}

func (d *ComputerDescriptorImpl) Colors() model.ColorPalette {
	return d.state().Colors
}

func (d *ComputerDescriptorImpl) ComputerLabel() string {
	return d.state().Label
}

func (d *ComputerDescriptorImpl) MessageBufferMap() model.BufferMap {
	return d.state().MessageBufferMap
}

func (d *ComputerDescriptorImpl) ConnectedSince() time.Time {
	return d.state().ConnectedSince
}

func (d *ComputerDescriptorImpl) Close() error {
	d.state().Activated = false

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
