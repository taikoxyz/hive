package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ethereum/hive/internal/libdocker"
	"github.com/ethereum/hive/internal/libhive"
	"github.com/stretchr/testify/assert"
)

type HiveConfig struct {
	ResultsRoot string
	Simulator   string
	Clients     string
}

type HiveFramework struct {
	*HiveConfig
	env        libhive.SimEnv
	simList    []string
	clientList []libhive.ClientDesignator
	runner     *libhive.Runner
}

func NewHiveFramework(config *HiveConfig) (*HiveFramework, error) {
	inv, err := libhive.LoadInventory(".")
	if err != nil {
		return nil, err
	}
	simList, err := inv.MatchSimulators(config.Simulator)
	if err != nil {
		return nil, err
	}

	builder, cb, err := libdocker.Connect("", &libdocker.Config{
		Inventory:           inv,
		PullEnabled:         false,
		UseCredentialHelper: false,
		ContainerOutput:     os.Stderr,
		BuildOutput:         os.Stderr,
	})
	if err != nil {
		return nil, err
	}
	clients, err := libhive.ParseClientList(&inv, config.Clients)
	if err != nil {
		return nil, err
	}

	return &HiveFramework{
		HiveConfig: config,
		env: libhive.SimEnv{
			LogDir:             config.ResultsRoot,
			SimLogLevel:        3,
			SimTestPattern:     "",
			SimParallelism:     1,
			SimRandomSeed:      0,
			SimDurationLimit:   0,
			ClientStartTimeout: 3 * time.Minute,
		},
		simList:    simList,
		clientList: clients,
		runner:     libhive.NewRunner(inv, builder, cb),
	}, nil
}

func (h *HiveFramework) Build(ctx context.Context) error {
	return h.runner.Build(ctx, h.clientList, h.simList)
}

func (h *HiveFramework) Run(ctx context.Context) (int, error) {
	var failCount int
	for _, sim := range h.simList {
		result, err := h.runner.Run(ctx, sim, h.env)
		if err != nil {
			return 0, err
		}
		failCount += result.TestsFailed
	}
	return failCount, nil
}

func TestTaikoSimulator(t *testing.T) {
	hiveTest, err := NewHiveFramework(&HiveConfig{
		ResultsRoot: "workspace/logs",
		Simulator:   "taiko",
		Clients:     "taiko/geth,taiko/taiko-geth,taiko/driver,taiko/proposer,taiko/prover",
	})
	assert.NoError(t, err)

	err = hiveTest.Build(context.Background())
	assert.NoError(t, err)

	failCount, err := hiveTest.Run(context.Background())
	assert.NoError(t, err)

	t.Log(failCount)
}
