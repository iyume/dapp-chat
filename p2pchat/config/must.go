package config

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func MustHexToECDSA(hexkey string) *ecdsa.PrivateKey {
	key, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		panic(err)
	}
	return key
}
