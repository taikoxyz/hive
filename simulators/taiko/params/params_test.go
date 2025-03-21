package params

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestENV(t *testing.T) {
	assert.Equal(t, true, len(EnvParams()) > 0)
	assert.Equal(t, true, len(ContractTxs) > 0)
	t.Logf(ContractTxs[0].Hash().String())
}
