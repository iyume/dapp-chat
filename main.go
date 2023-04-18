package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	chatABI "github.com/iyume/go-blockchain-chat/go-chat"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const key = `data/keystore/UTC--2023-04-18T12-37-51.842961176Z--8cedb7c6af8a7781ec89bb84900768de99c8235b`

func main() {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial("data/geth.ipc")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	file, err := os.ReadFile(key)
	if err != nil {
		log.Fatalf("Failed to read key: %v", err)
	}
	auth, err := bind.NewTransactorWithChainID(strings.NewReader(string(file)), "123", big.NewInt(12345))
	auth.GasLimit = 8000000
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	// Deploy the contract passing the newly created `auth` and `conn` vars
	address, tx, instance, err := chatABI.DeployChat(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy new storage contract: %v", err)
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())

	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	// function call
	res, err := instance.AddTwo(&bind.CallOpts{Pending: true}, big.NewInt(3), big.NewInt((5)))
	if err != nil {
		log.Fatalf("Failed to retrieve result: %v", err)
	}
	fmt.Println("Result:", res)
}
