package main

import (
	"crypto/ecdsa"
	"log"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

// for test
var bootnode = ""

func getPrivateKeyJson(file string) *ecdsa.PrivateKey {
	keyjson, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln("Error load keyfile:", err)
	}
	key, err := keystore.DecryptKey(keyjson, "123")
	if err != nil {
		log.Fatalln(err)
	}
	return key.PrivateKey
}

func getGenerateKey() *ecdsa.PrivateKey {
	priv, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalln(err)
	}
	return priv
}

func main() {
	SetLogLevel(LvlFromString("debug"))

	if err := p2pserver.Start(); err != nil {
		log.Fatalln(err)
	}
	defer p2pserver.Stop()
	// srv.LocalNode().Node() ensure localnode exists. srv.Self() will create it.
	log.Println("Started P2P networking at", p2pserver.LocalNode().Node().URLv4())

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Started http server at", "http://"+listener.Addr().String())
	if err := server.Serve(listener); err != nil {
		log.Println(err)
	}
}
