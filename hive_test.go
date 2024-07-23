package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaikoSimulator(t *testing.T) {
	hiveTest, err := NewHiveFramework(&HiveConfig{
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

	t.Log(failCount)
}
