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
const keyfile = `data/keystore/UTC--2023-04-18T12-49-23.162529799Z--d542be4551d114a7a2b544bafb7a9feba8784e69`
const passphrase = "123"
const chainID = 12345
const contractAddress = "0x7b17f3af835b9ef46e1bad6f0d3f8352882f672f"

func main() {
	conn, err := ethclient.Dial(ipc_path)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	instance, err := chatABI.NewChat(common.HexToAddress(
		contractAddress), conn)
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
	auth.GasPrice = big.NewInt(1)
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
