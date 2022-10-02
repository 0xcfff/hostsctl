package hosts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBackend(t *testing.T) {
	t.Run("create default backend", func(t *testing.T) {
		backend := NewBackend(nil)
		state, err := backend.ReadState()
		assert.NoError(t, err)
		assert.NotNil(t, state)
	})
}
