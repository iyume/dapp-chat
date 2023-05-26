package main

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	log.Printf("%x\n", crypto.FromECDSA(key))
}
