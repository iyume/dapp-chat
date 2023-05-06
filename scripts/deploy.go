package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	chatABI "github.com/iyume/dapp-chat/go-dapp-chat"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const chainID = 12345
const ipc_path = "data/geth.ipc"
const keyfile = `data/keystore/UTC--2023-04-18T12-49-23.162529799Z--d542be4551d114a7a2b544bafb7a9feba8784e69`
const passphrase = "123"
const deployTimeout = 5

func main() {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(ipc_path)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
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
	// NOTE: transaction require at least one gasprice (miner limits)
	auth.GasPrice = big.NewInt(1)
	address, tx, _, err := chatABI.DeployChat(auth, conn)
	if err != nil {
		log.Fatalln("Failed to deploy contract:", err)
	}
	// Use format because `address.String()` has mixed case
	fmt.Printf("Contract pending address: 0x%x\n", address)
	fmt.Println("Transaction waiting to be verified:", tx.Hash())
}
