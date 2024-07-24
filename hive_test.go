package main

import (
	"context"
	"github.com/ethereum/hive/taiko"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaikoSimulator(t *testing.T) {
	hiveTest, err := taiko.NewHiveFramework(&taiko.HiveConfig{
		SimPattern: "taiko/driver",
		Clients: []string{
			"taiko/geth",
			"taiko/taiko-geth",
			"taiko/driver",
			"taiko/proposer",
		},
	})
	assert.NoError(t, err)

	failCount, err := hiveTest.Run(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, 0, failCount)
}
