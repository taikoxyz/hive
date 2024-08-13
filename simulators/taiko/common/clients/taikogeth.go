package clients

import (
	"fmt"
)

type TaikoGethClient struct {
	*BaseNode
}

func (t *TaikoGethClient) EngineURL() string {
	return fmt.Sprintf("http://%v:%d", t.NetworkIP(), EthEngineRPC)
}
