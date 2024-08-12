package utils

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ExtraVanity = 32                     // Fixed number of extra-data prefix bytes reserved for signer vanity
	ExtraSeal   = crypto.SignatureLength // Fixed number of extra-data suffix bytes reserved for signer seal
)

type TestAccount struct {
	key     *ecdsa.PrivateKey
	address *common.Address
	index   uint64
}

func (a *TestAccount) GetKey() *ecdsa.PrivateKey {
	return a.key
}

func (a *TestAccount) GetAddress() common.Address {
	if a.address == nil {
		key := a.key
		addr := crypto.PubkeyToAddress(key.PublicKey)
		a.address = &addr
	}
	return *a.address
}

func (a *TestAccount) GetIndex() uint64 {
	return a.index
}
