package preconf

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"taiko/params"
	"testing"
)

func TestNewProposer(t *testing.T) {
	propose, err := NewProposer()
	assert.Nil(t, err)
	assert.NotNil(t, propose)

	l2Cli, err := ethclient.Dial(params.ParamByKey("L2_HTTP"))
	assert.Nil(t, err)

	err = ProposeTxLists(context.Background(), propose, l2Cli)
	assert.Nil(t, err)
}
