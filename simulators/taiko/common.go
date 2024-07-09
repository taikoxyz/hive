package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/hive/hivesim"
	"github.com/shogo82148/go-tap"
	"io"
	"os"
	"os/exec"
)

const (
	envFile = "/taiko/.env"
	network = "hive_taiko_network"
)

var networkCreated = make(map[hivesim.SuiteID]bool)

// createNetwork ensures there is a separate network to be able to send the client traffic
// from two separate IP addrs.
func createAndConnectNetwork(t *hivesim.T, container string) {
	if !networkCreated[t.SuiteID] {
		if err := t.Sim.CreateNetwork(t.SuiteID, network); err != nil {
			t.Fatal("can't create network:", err)
		}
		if err := t.Sim.ConnectContainer(t.SuiteID, network, "simulation"); err != nil {
			t.Fatal("can't connect simulation to network:", err)
		}
		networkCreated[t.SuiteID] = true
	}

	if err := t.Sim.ConnectContainer(t.SuiteID, network, container); err != nil {
		t.Fatal("can't connect container to network:", err)
	}
}

func runTAP(t *hivesim.T, clientName string, cmd *exec.Cmd) error {
	// Set up output streams.
	cmd.Stderr = os.Stderr
	output, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("can't set up test command stdout pipe: %v", err)
	}
	defer output.Close()

	// Forward TAP output to the simulator log.
	outputTee := io.TeeReader(output, os.Stdout)

	// Run the test command.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("can't start test command: %v", err)
	}
	if err := reportTAP(t, clientName, outputTee); err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return err
	}
	return cmd.Wait()
}

func reportTAP(t *hivesim.T, clientName string, output io.Reader) error {
	// Parse the output.
	parser, err := tap.NewParser(output)
	if err != nil {
		return fmt.Errorf("error parsing TAP: %v", err)
	}
	for {
		test, err := parser.Next()
		if test == nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		// Forward result to hive.
		name := fmt.Sprintf("%s (%s)", test.Description, clientName)
		testID, err := t.Sim.StartTest(t.SuiteID, hivesim.TestStartInfo{Name: name})
		if err != nil {
			return fmt.Errorf("can't report sub-test result: %v", err)
		}
		result := hivesim.TestResult{Pass: test.Ok, Details: test.Diagnostic}
		t.Sim.EndTest(t.SuiteID, testID, result)
	}
}

func testGeth(t *hivesim.T, c *hivesim.Client) {
	url := fmt.Sprintf("http://%v:8545", c.IP)

	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ChainID(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("container id: %s", c.Container)
}
