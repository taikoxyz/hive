package clients

import (
	"fmt"
	"time"

	"github.com/ethereum/hive/hivesim"
)

type HiveOptionsGenerator func() ([]hivesim.StartOption, error)

type HiveManagedClient struct {
	*hivesim.T
	OptionsGenerator     HiveOptionsGenerator
	HiveClientDefinition *hivesim.ClientDefinition

	simClient         *hivesim.Client
	extraStartOptions []hivesim.StartOption

	Network   string
	networkIP string
}

func (h *HiveManagedClient) IsRunning() bool {
	return h.simClient != nil
}

func (h *HiveManagedClient) Start() error {
	h.Logf("Starting %s client", h.ClientType())
	opts, err := h.OptionsGenerator()
	if err != nil {
		return fmt.Errorf("unable to get start options: %v", err)
	}

	if opts == nil {
		opts = make([]hivesim.StartOption, 0)
	}

	if h.extraStartOptions != nil {
		opts = append(opts, h.extraStartOptions...)
	}

	h.simClient = h.StartClient(h.HiveClientDefinition.Name, opts...)
	if h.simClient == nil {
		return fmt.Errorf("unable to launch client")
	}
	h.Logf(
		"Started client %s, container %s",
		h.ClientType(),
		h.simClient.Container,
	)
	return nil
}

func (h *HiveManagedClient) PauseClient() {
	h.Logf("Pausing %s client", h.ClientType())
	err := h.Sim.PauseClient(h.SuiteID, h.TestID, h.simClient.Container)
	if err != nil {
		h.Fatalf("Error pausing %s client: %v", h.ClientType(), err)
	}
	time.Sleep(time.Second)
}

func (h *HiveManagedClient) UnpauseClient() {
	h.Logf("Unpausing %s client", h.ClientType())
	err := h.Sim.UnpauseClient(h.SuiteID, h.TestID, h.simClient.Container)
	if err != nil {
		h.Fatalf("Error unpausing %s client: %v", h.ClientType(), err)
	}
}

func (h *HiveManagedClient) NetworkIP() string {
	if h.networkIP != "" {
		return h.networkIP
	}

	var err error
	h.networkIP, err = h.Sim.ContainerNetworkIP(h.SuiteID, h.Network, h.simClient.Container)
	if err != nil {
		// h.Logf("Error getting network IP: %v", err)
		return h.GetHost()
	}
	h.Logf("%s network IP %v", h.ClientType(), h.networkIP)

	return h.networkIP
}

func (h *HiveManagedClient) GetHost() string {
	if h.simClient == nil {
		return ""
	}
	return h.simClient.IP.String()
}

func (h *HiveManagedClient) Shutdown() error {
	h.Logf("Shutting down %s client, container %s", h.ClientType(), h.simClient.Container)
	if err := h.Sim.StopClient(h.SuiteID, h.TestID, h.simClient.Container); err != nil {
		return err
	}
	h.simClient = nil
	return nil
}

func (h *HiveManagedClient) GetEnodeURL() (string, error) {
	return h.simClient.EnodeURL()
}

func (h *HiveManagedClient) ClientType() string {
	return h.HiveClientDefinition.Name
}
