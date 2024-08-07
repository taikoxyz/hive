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
	T                    *hivesim.T
	OptionsGenerator     HiveOptionsGenerator
	HiveClientDefinition *hivesim.ClientDefinition
	Port                 int64

	Client            *hivesim.Client
	extraStartOptions []hivesim.StartOption
}

func (h *HiveManagedClient) HiveClient() *hivesim.Client {
	return h.Client
}

func (h *HiveManagedClient) IsRunning() bool {
	return h.Client != nil
}

func (h *HiveManagedClient) Start() error {
	h.T.Logf("Starting client %s", h.ClientType())
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

	h.Client = h.T.StartClient(h.HiveClientDefinition.Name, opts...)
	if h.Client == nil {
		return fmt.Errorf("unable to launch client")
	}
	h.T.Logf(
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
	if err := h.T.Sim.StopClient(h.T.SuiteID, h.T.TestID, h.Client.Container); err != nil {
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
