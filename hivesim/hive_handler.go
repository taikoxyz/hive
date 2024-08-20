package hivesim

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/hive/internal/libdocker"
	"github.com/ethereum/hive/internal/libhive"
)

var (
	DefaultHiveConfig = &HiveConfig{
		ResultsRoot:    "workspace/logs",
		Loglevel:       3,
		SimParallelism: 1,
		SimLogLevel:    4,
		ClientTimeOut:  time.Minute * 3,
		BaseDir:        ".",
		DockerPull:     true,
	}
)

func hiveConfigWithDefault(cfg *HiveConfig) *HiveConfig {
	if cfg.ResultsRoot == "" {
		cfg.ResultsRoot = DefaultHiveConfig.ResultsRoot
	}
	if cfg.Loglevel == 0 {
		cfg.Loglevel = DefaultHiveConfig.Loglevel
	}
	if cfg.SimParallelism == 0 {
		cfg.SimParallelism = DefaultHiveConfig.SimParallelism
	}
	if cfg.SimLogLevel == 0 {
		cfg.SimLogLevel = DefaultHiveConfig.SimLogLevel
	}
	if cfg.ClientTimeOut == 0 {
		cfg.ClientTimeOut = DefaultHiveConfig.ClientTimeOut
	}
	if cfg.BaseDir == "" {
		cfg.BaseDir = DefaultHiveConfig.BaseDir
	}
	if cfg.DockerPull == false {
		cfg.DockerPull = DefaultHiveConfig.DockerPull
	}
	return cfg
}

type HiveConfig struct {
	ResultsRoot string
	Loglevel    int

	DockerPull      bool
	DockerNoCache   string
	BuildOutput     bool
	ContainerOutput bool

	SimPattern     string
	SimTestPattern string
	SimParallelism int
	SimRandomSeed  int
	SimTimeLimit   time.Duration
	SimLogLevel    int

	Clients       []string
	ClientTimeOut time.Duration

	BaseDir string
}

type HiveFramework struct {
	*HiveConfig
	env        libhive.SimEnv
	simList    []string
	clientList []libhive.ClientDesignator
	runner     *libhive.Runner
}

func NewHiveFramework(config *HiveConfig) (*HiveFramework, error) {
	cfg := hiveConfigWithDefault(config)
	inv, err := libhive.LoadInventory(config.BaseDir)
	if err != nil {
		return nil, err
	}
	simList, err := inv.MatchSimulators(cfg.SimPattern)
	if err != nil {
		return nil, err
	}

	dockerConfig := &libdocker.Config{
		Inventory:           inv,
		PullEnabled:         config.DockerPull,
		UseCredentialHelper: false,
	}
	if cfg.DockerNoCache != "" {
		re, err := regexp.Compile(cfg.DockerNoCache)
		if err != nil {
			fmt.Fprintln(os.Stderr, "bad --docker-nocache regular expression:", err)
			os.Exit(1)
		}
		dockerConfig.NoCachePattern = re
	}
	if cfg.ContainerOutput {
		dockerConfig.ContainerOutput = os.Stderr
	}
	if cfg.BuildOutput {
		dockerConfig.BuildOutput = os.Stderr
	}

	builder, cb, err := libdocker.Connect("", dockerConfig)
	if err != nil {
		return nil, err
	}
	clients, err := libhive.ParseClientList(&inv, strings.Join(cfg.Clients, ","))
	if err != nil {
		return nil, err
	}

	return &HiveFramework{
		HiveConfig: cfg,
		env: libhive.SimEnv{
			LogDir:             cfg.ResultsRoot,
			SimLogLevel:        cfg.SimLogLevel,
			SimTestPattern:     cfg.SimTestPattern,
			SimParallelism:     cfg.SimParallelism,
			SimRandomSeed:      cfg.SimRandomSeed,
			SimDurationLimit:   cfg.SimTimeLimit,
			ClientStartTimeout: cfg.ClientTimeOut,
		},
		simList:    simList,
		clientList: clients,
		runner:     libhive.NewRunner(inv, builder, cb),
	}, nil
}

func (h *HiveFramework) build(ctx context.Context) error {
	return h.runner.Build(ctx, h.clientList, h.simList)
}

func (h *HiveFramework) Run(ctx context.Context) (int, error) {
	// build simulators and clients.
	if err := h.build(ctx); err != nil {
		return 0, err
	}

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
