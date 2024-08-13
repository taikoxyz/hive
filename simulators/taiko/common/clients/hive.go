package clients

import (
	"fmt"
	"net"

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

	Client            *hivesim.Client
	extraStartOptions []hivesim.StartOption

	Network   string
	networkIP string
}

func (h *HiveManagedClient) HiveClient() *hivesim.Client {
	return h.Client
}

func (h *HiveManagedClient) IsRunning() bool {
	return h.Client != nil
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

	h.Client = h.StartClient(h.HiveClientDefinition.Name, opts...)
	if h.Client == nil {
		return fmt.Errorf("unable to launch client")
	}
	h.Logf(
		"Started client %s, container %s",
		h.ClientType(),
		h.Client.Container,
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

func (h *HiveManagedClient) NetworkIP() string {
	if h.networkIP != "" {
		return h.networkIP
	}

	var err error
	h.networkIP, err = h.Sim.ContainerNetworkIP(h.SuiteID, h.Network, h.Client.Container)
	if err != nil {
		h.Logf("Error getting network IP: %v", err)
		return h.GetHost()
	}
	h.Logf("taiko geth network IP %v", h.networkIP)

	return h.networkIP
}

func (h *HiveManagedClient) GetAddress() string {
	if h.Client == nil {
		return ""
	}
	if h.Port > 0 {
		return fmt.Sprintf("http://%s:%d", h.Client.IP, h.Port)
	}
	return fmt.Sprintf("http://%s", h.Client.IP)
}

func (h *HiveManagedClient) GetIP() net.IP {
	if h.Client == nil {
		return net.IP{}
	}
	return h.Client.IP
}

func (h *HiveManagedClient) GetHost() string {
	if h.Client == nil {
		return ""
	}
	return h.Client.IP.String()
}

func (h *HiveManagedClient) Shutdown() error {
	h.Logf("Shutting down %s client, container %s", h.ClientType(), h.Client.Container)
	if err := h.Sim.StopClient(h.SuiteID, h.TestID, h.Client.Container); err != nil {
		return err
	}
	h.Client = nil
	return nil
}

func (h *HiveManagedClient) GetEnodeURL() (string, error) {
	return h.Client.EnodeURL()
}

func (h *HiveManagedClient) ClientType() string {
	return h.HiveClientDefinition.Name
}
