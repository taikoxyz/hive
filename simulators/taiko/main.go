package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"os"
	"os/exec"

	"github.com/ethereum/hive/hivesim"
	"github.com/shogo82148/go-tap"
)

func main() {
	init := hivesim.Suite{
		Name:        "init",
		Description: `Test framework initialization`,
	}
	init.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "contract",
		Description: "Deploy taiko contract on l1 chain",
		Run:         deployL1Contract,
		AlwaysRun:   false,
	})
	init.Add(hivesim.ClientTestSpec{
		Role:        "geth",
		Name:        "contract",
		Description: "Show contract addresses",
		Run:         showEnv,
		AlwaysRun:   false,
	})
	hivesim.MustRun(hivesim.New(), init)
}

// Deploy a contract on L1
func deployL1Contract(t *hivesim.T, c *hivesim.Client) {
	url := fmt.Sprintf("http://%v:8545", c.IP)
	if err := os.Setenv("L1_NODE_HTTP_ENDPOINT", url); err != nil {
		t.Fatal(err)
	}
	// deploy l1 contract.
	cmd := exec.Command("sh", "/taiko/deploy_l1_contract.sh")
	if err := runTAP(t, c.Type, cmd); err != nil {
		t.Fatal(err)
	}
}

func showEnv(t *hivesim.T, c *hivesim.Client) {
	if err := godotenv.Load("/taiko/.env"); err != nil {
		t.Fatal(err)
	}
	l1Address := os.Getenv("TAIKO_L1_ADDRESS")
	fmt.Println("l1Address: ", l1Address)
	tokenAddress := os.Getenv("TAIKO_TOKEN_ADDRESS")
	fmt.Println("tokenAddress: ", tokenAddress)
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
