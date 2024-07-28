package suite_reorg

import (
	"bytes"
	"github.com/ethereum/hive/hivesim"
	"github.com/ethereum/hive/simulators/eth2/common/clients"
	"github.com/ethereum/hive/simulators/eth2/common/testnet"
	"github.com/protolambda/ztyp/codec"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCC(t *testing.T) {
	var clientTypes = []*hivesim.ClientDefinition{
		{
			Name:    "go-ethereum",
			Version: "client-1-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"eth1"}},
		},
		{
			Name:    "prysm-bn",
			Version: "client-2-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"beacon"}},
		}, {
			Name:    "prysm-vc",
			Version: "client-3-version",
			Meta:    hivesim.ClientMetadata{Roles: []string{"validator"}},
		},
	}

	clientsByRole := clients.ClientsByRole(clientTypes)

	clientCombinations := clientsByRole.Combinations()

	mnemonic := "couple kiwi radio river setup fortune hunt grief buddy forward perfect empty slim wear bounce drift execute nation tobacco dutch chapter festival ice fog"

	for _, test := range Tests {
		keys := test.GetValidatorKeys(mnemonic)
		env := &testnet.Environment{
			Clients:    clientsByRole,
			Validators: keys,
		}
		config := test.GetTestnetConfig(clientCombinations)

		prep, err := testnet.PrepareTestnet(env, config)
		assert.NoError(t, err)

		var stateBytes bytes.Buffer

		err = prep.BeaconGenesis.Serialize(codec.NewEncodingWriter(&stateBytes))
		assert.NoError(t, err)

		err = os.WriteFile("/tmp/genesis.ssz", stateBytes.Bytes(), 0644)
		assert.NoError(t, err)

		t.Log(prep)
	}
}
