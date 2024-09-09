package clients

import (
	"fmt"
	"net"
	"time"

	"github.com/ethereum/hive/hivesim"
	"github.com/marioevz/eth-clients/clients"
)

var _ clients.ManagedClient = &HiveManagedClient{}

type HiveOptionsGenerator func() ([]hivesim.StartOption, error)

type HiveManagedClient struct {
	*hivesim.T
	OptionsGenerator     HiveOptionsGenerator
	HiveClientDefinition *hivesim.ClientDefinition
	Port                 int64

	simClient         *hivesim.Client
	extraStartOptions []hivesim.StartOption

	Network   string
	networkIP string
}

func (h *HiveManagedClient) HiveClient() *hivesim.Client {
	return h.simClient
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

func (h *HiveManagedClient) AddStartOption(opts ...interface{}) {
	if h.extraStartOptions == nil {
		h.extraStartOptions = make([]hivesim.StartOption, 0)
	}
	for _, o := range opts {
		if o, ok := o.(hivesim.StartOption); ok {
			h.extraStartOptions = append(h.extraStartOptions, o)
		}
	}
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
		h.Logf("Error getting network IP: %v", err)
		return h.GetHost()
	}
	h.Logf("%s network IP %v", h.ClientType(), h.networkIP)

	return h.networkIP
}

func (h *HiveManagedClient) GetAddress() string {
	if h.simClient == nil {
		return ""
	}
	if h.Port > 0 {
		return fmt.Sprintf("http://%s:%d", h.simClient.IP, h.Port)
	}
	return fmt.Sprintf("http://%s", h.simClient.IP)
}

func (h *HiveManagedClient) GetIP() net.IP {
	if h.simClient == nil {
		return net.IP{}
	}
	return h.simClient.IP
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
