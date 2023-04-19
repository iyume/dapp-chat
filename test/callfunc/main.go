package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	chatABI "github.com/iyume/go-blockchain-chat/go-chat"
)

const ipc_path = "data/geth.ipc"
const keyfile = `data/keystore/UTC--2023-04-18T12-37-51.842961176Z--8cedb7c6af8a7781ec89bb84900768de99c8235b`
const passphrase = "123"
const chainID = 12345

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
	keyfilebuf, err := os.ReadFile(keyfile)
	if err != nil {
		log.Fatalf("Failed to read keyfile: %v", err)
	}
	auth, err := bind.NewTransactorWithChainID(
		strings.NewReader(string(keyfilebuf)), passphrase, big.NewInt(chainID),
	)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	session := chatABI.ChatSession{Contract: instance,
		CallOpts: bind.CallOpts{Pending: true}, TransactOpts: *auth}
	res, err := session.GetValue()
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}
	fmt.Printf("Call result: %v\n", res)
	// transaction from auth to contract
	_, err = session.Increase(big.NewInt(12))
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}
	res, err = session.GetValue()
	if err != nil {
		log.Fatalf("Failed to call function: %v", err)
	}
	fmt.Printf("Call result: %v\n", res)
}
