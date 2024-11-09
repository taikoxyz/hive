package clients

import (
	"fmt"
)

const PortBeaconTCP = 9000

type BeaconClient struct {
	*HiveManagedClient
}

func (bn *BeaconClient) BeaconURL() string {
	return fmt.Sprintf("http://%s:%d", bn.NetworkIP(), BeaconPort)
}

func (bn *BeaconClient) ClientName() string {
	name := bn.ClientType()
	if len(name) > 3 && name[len(name)-3:] == "-bn" {
		name = name[:len(name)-3]
	}
	return name
}

type BlockV2OptimisticResponse struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
}
