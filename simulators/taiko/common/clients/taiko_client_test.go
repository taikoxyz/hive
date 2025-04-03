package clients

import (
	"github.com/stretchr/testify/assert"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/cmd/flags"
	"testing"
)

func TestMockDriver(t *testing.T) {
	driverClient := &MockDriver{}
	err := NewTaikoClient(driverClient, flags.DriverFlags)
	assert.NoError(t, err)
}
