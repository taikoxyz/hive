package preconf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProposer(t *testing.T) {
	propose, err := NewProposer()
	assert.Nil(t, err)
	assert.NotNil(t, propose)
}
