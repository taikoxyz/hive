package main

import (
	"fmt"
	"github.com/ethereum/hive/hivesim"
	"github.com/shogo82148/go-tap"
	"io"
	"os"
	"os/exec"
)

func main() {
	suite := hivesim.Suite{
		Name:        "taiko",
		Description: ``,
	}
	suite.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "deploy l1 contract",
		Description: "deploy l1 contract",
		Run:         deployL1Contract,
		AlwaysRun:   false,
	})
}

// Deploy a contract on L1
func deployL1Contract(t *hivesim.T, c *hivesim.Client) {
	url := fmt.Sprintf("http://%v:8545", c.IP)
	if err := os.Setenv("L2_EXECUTION_ENGINE_HTTP_ENDPOINT", url); err != nil {
		t.Fatal(err)
	}

	// deploy l1 contract.
	cmd := exec.Command("sh", "/taiko/deploy_l1_contract.sh")
	if err := runTAP(t, c.Type, cmd); err != nil {
		t.Fatal(err)
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
	return nil
}
