package main

import (
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"time"
)

type containerParams struct {
	// This filters client types by role.
	// If no role is specified, the test runs for all available client types.
	Role string

	// Parameters and Files are launch options for client instances.
	Parameters hivesim.Params
	Files      map[string]string
}

type TaikoTestSpec struct {
	cParams    map[string]*containerParams
	containers map[string]*hivesim.Client

	suites []hivesim.Suite
	specs  []*hivesim.ClientTestSpec
}

func NewCommonSpec() *TaikoTestSpec {
	return &TaikoTestSpec{
		cParams:    make(map[string]*containerParams),
		containers: make(map[string]*hivesim.Client),
	}
}

func (c *TaikoTestSpec) StartTest() {
	suite := hivesim.Suite{
		Name: "taiko-test-enter",
		Tests: []hivesim.AnyTest{
			hivesim.TestSpec{
				Name: "taiko-test-enter",
				Run:  c.tests,
			},
		},
	}
	hivesim.MustRun(hivesim.New(), suite)
}

func (c *TaikoTestSpec) tests(t *hivesim.T) {
	cNames, err := t.Sim.ClientTypes()
	if err != nil {
		t.Fatalf("failed to get client types: %v", err)
	}
	for _, ct := range c.cParams {
		for _, cn := range cNames {
			if ct.Role != "" && !cn.HasRole(ct.Role) {
				continue
			}
			client := t.StartClient(cn.Name, ct.Parameters, hivesim.WithStaticFiles(ct.Files))
			c.containers[ct.Role] = client
			t.Logf("started container %s: %s", ct.Role, client.Container)
		}
	}
	defer func() {
		for role, cn := range c.containers {
			if err := t.Sim.StopClient(t.SuiteID, t.TestID, cn.Container); err != nil {
				t.Errorf("failed to stop container %s: %v", role, err)
			}
		}
	}()

	for _, suite := range c.suites {
		for i, test := range suite.Tests {
			switch tp := test.(type) {
			case hivesim.ClientTestSpec:
				if _, ok := c.containers[tp.Role]; !ok {
					t.Errorf("container %s does not exist", tp.Role)
					continue
				}
				suite.Tests[i] = hivesim.TestSpec{
					Name:        tp.Name,
					Description: tp.Description,
					Category:    tp.Category,
					AlwaysRun:   tp.AlwaysRun,
					Run: func(t *hivesim.T) {
						tp.Run(t, c.containers[tp.Role])
					},
				}
			case hivesim.TestSpec:
			}
		}
	}

	// Run all the suites
	hivesim.MustRun(hivesim.New(), c.suites...)

	time.Sleep(100 * time.Second)

	t.Log("all the tests done")
}

func (c *TaikoTestSpec) AddSuite(suites ...hivesim.Suite) error {
	for _, suite := range suites {
		for _, test := range suite.Tests {
			switch tp := test.(type) {
			case hivesim.ClientTestSpec:
				if _, ok := c.cParams[tp.Role]; ok {
					return fmt.Errorf("container role %s already exist", tp.Role)
				}
				c.cParams[tp.Role] = &containerParams{
					Role:       tp.Role,
					Parameters: tp.Parameters,
					Files:      tp.Files,
				}
			case hivesim.TestSpec:
			default:
				return fmt.Errorf("unsupported type: %v", tp)
			}
		}
		c.suites = append(c.suites, suite)
	}

	return nil
}
