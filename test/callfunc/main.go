package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	chatABI "github.com/iyume/go-blockchain-chat/go-chat"
)

const ipc_path = "data/geth.ipc"

func main() {
	conn, err := ethclient.Dial(ipc_path)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	instance, err := chatABI.NewChat(common.HexToAddress(
		"0x695483a11e353632ed20780a33755e2f64a0361e"), conn)
	if err != nil {
		log.Fatalf("Failed to create contract instance: %v", err)
	}
	session := chatABI.ChatSession{Contract: instance,
		CallOpts: bind.CallOpts{Pending: true}}
	res, err := session.AddTwo(big.NewInt(12), big.NewInt(13))
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}
	fmt.Printf("Call result: %v", res)
}
