package types

import (
	"cc-rshell-server/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/olahol/melody.v1"
	"testing"
	"time"
)

func TestComputerDescriptorImpl_Init(t *testing.T) {
	t.Parallel()

	d := &ComputerDescriptorImpl{
		Session: &melody.Session{
			Keys: map[string]interface{}{},
		},
	}

	d.Init()

	assert.NotZero(t, d.ConnectedSince())
	assert.False(t, d.Activated())
}

func TestComputerDescriptorImpl_Activate(t *testing.T) {
	t.Parallel()

	d := &ComputerDescriptorImpl{
		Session: &melody.Session{
			Keys: map[string]interface{}{},
		},
	}

	d.Init()
	d.Activate(42, "nasapc", map[string]interface{}{
		"enter": 257,
	}, nil)

	assert.NotZero(t, d.ConnectedSince())
	assert.True(t, d.Activated())
	assert.Equal(t, model.ComputerID(42), d.ComputerID())
	assert.Equal(t, "nasapc", d.ComputerLabel())
	assert.Equal(t, 257, d.KeyCodes()["enter"])
}

type ComputerDescriptorImplMock struct {
	ComputerDescriptorImpl
	stopped chan struct{}
}

func (m *ComputerDescriptorImplMock) Close() error {
	close(m.stopped)
	return nil
}

func TestComputerDescriptorImpl_ActivationTimeout(t *testing.T) {
	t.Parallel()

	d := ComputerDescriptorImpl{
		Session: &melody.Session{
			Keys: map[string]interface{}{},
		},
	}

	stopChan := make(chan struct{})
	m := ComputerDescriptorImplMock{
		ComputerDescriptorImpl: d,
		stopped:                stopChan,
	}

	m.initTimeout(0, m.Close)
	select {
	case <-stopChan:
		// stopped after timeout
	case <-time.After(1 * time.Second):
		assert.FailNow(t, "not stopped after timeout")
	}
}
