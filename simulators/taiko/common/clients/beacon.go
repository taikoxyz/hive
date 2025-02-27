package clients

import (
	"fmt"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/taikoxyz/taiko-mono/packages/taiko-client/pkg/rpc"
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

func MakeBlobs(txListBytes []byte) ([]*eth.Blob, error) {
	var blobs []*eth.Blob
	for start := 0; start < len(txListBytes); start += rpc.BlobBytes {
		end := start + rpc.BlobBytes
		if end > len(txListBytes) {
			end = len(txListBytes)
		}

		var blob = &eth.Blob{}
		if err := blob.FromData(txListBytes[start:end]); err != nil {
			return nil, err
		}
		blob.KZGBlob()

		blobs = append(blobs, blob)
	}

	return blobs, nil
}
