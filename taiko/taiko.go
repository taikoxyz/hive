package taiko

import (
	"fmt"
	"sync"

	"github.com/ethereum/hive/hivesim"
)

type TaikoTestSpec struct {
	*hivesim.T

	cDefines   []*hivesim.ClientDefinition
	containers map[string]*hivesim.Client

	wait   sync.WaitGroup
	stopCh chan struct{}
}

func NewCommonSpec() *TaikoTestSpec {
	taiko := &TaikoTestSpec{
		containers: make(map[string]*hivesim.Client),
		stopCh:     make(chan struct{}),
	}

	taiko.init()

	return taiko
}

func (t *TaikoTestSpec) Release() {
	// Wait until all the suites finished running.
	t.wait.Wait()

	// Release containers.
	for role, cn := range t.containers {
		if err := t.Sim.StopClient(t.SuiteID, t.TestID, cn.Container); err != nil {
			t.Errorf("failed to stop container %s: %v", role, err)
		}
	}

	close(t.stopCh)
}

func (t *TaikoTestSpec) init() {
	t.wait.Add(1)

	suite := hivesim.Suite{
		Name: "taiko-test-spec",
		Tests: []hivesim.AnyTest{
			hivesim.TestSpec{
				Name: "taiko-test-spec",
				Run:  t.run,
			},
		},
	}
	go hivesim.MustRun(hivesim.New(), suite)

	t.wait.Wait()
}

func (t *TaikoTestSpec) run(sim *hivesim.T) {
	t.T = sim

	defines, err := t.Sim.ClientTypes()
	if err != nil {
		t.Fatalf("failed to get client types: %v", err)
	}

	// Show client types.
	names := make([]string, 0, len(defines))
	for _, cn := range defines {
		names = append(names, cn.Name)
	}
	t.Logf("docker image types: %v", names)

	t.cDefines = defines

	t.wait.Done()

	<-t.stopCh
	t.Logf("The taiko-test-spec is stopped.")
}

func (t *TaikoTestSpec) startClient(role string, params hivesim.Params, files map[string]string) (*hivesim.Client, error) {
	if client, ok := t.containers[role]; ok {
		t.Logf("container %s already started: %s", role, client.Container)
		return client, nil
	}
	for _, cn := range t.cDefines {
		if role != "" && !cn.HasRole(role) {
			continue
		}
		client := t.StartClient(cn.Name, params, hivesim.WithStaticFiles(files))
		t.containers[role] = client
		t.Logf("started container %s: %s", role, client.Container)
		return client, nil
	}
	return nil, fmt.Errorf("container role %s not match", role)
}

func (t *TaikoTestSpec) RunSuite(suites ...hivesim.Suite) {
	t.wait.Add(1)
	defer t.wait.Done()

	for _, suite := range suites {
		for i, test := range suite.Tests {
			switch tp := test.(type) {
			case hivesim.ClientTestSpec:
				// Init container.
				container, err := t.startClient(tp.Role, tp.Parameters, tp.Files)
				if err != nil {
					t.Errorf("failed to start container %s: %v", tp.Role, err)
					return
				}
				suite.Tests[i] = hivesim.TestSpec{
					Name:        tp.Name,
					Description: tp.Description,
					Category:    tp.Category,
					AlwaysRun:   tp.AlwaysRun,
					Run: func(t *hivesim.T) {
						tp.Run(t, container)
					},
				}
			case hivesim.TestSpec:
			default:
				t.Errorf("unsupported type: %v", tp)
			}
		}

		t.Logf("running suite %s", suite.Name)
		hivesim.MustRunSuite(hivesim.New(), suite)
	}
}
