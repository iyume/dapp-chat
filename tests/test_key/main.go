package main

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

func main() {
	key, _ := crypto.GenerateKey()
	log.Printf("%x\n", crypto.Keccak256(crypto.FromECDSAPub(&key.PublicKey)[1:]))
	// Convert to [32]byte because enode.ID has String()
	log.Printf("%x\n", enode.PubkeyToIDV4(&key.PublicKey).Bytes())
}
