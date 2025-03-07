package preconf

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func parsePriv(data string) (*crypto.Secp256k1PrivateKey, error) {
	if len(data) > 2 && data[:2] == "0x" {
		data = data[2:]
	}
	b, err := hex.DecodeString(data)
	if err != nil {
		return nil, errors.New("p2p priv key is not formatted in hex chars")
	}
	p, err := crypto.UnmarshalSecp256k1PrivateKey(b)
	if err != nil {
		// avoid logging the priv key in the error, but hint at likely input length problem
		return nil, fmt.Errorf("failed to parse priv key from %d bytes", len(b))
	}
	return (p).(*crypto.Secp256k1PrivateKey), nil
}
