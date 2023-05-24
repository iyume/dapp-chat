package main

import (
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/server"
)

func main() {
	SetLogLevel(LvlFromString("debug"))

	bconfig := api.DefaultBackendConfig
	bconfig.Key, _ = crypto.GenerateKey()
	backend := api.NewBackend(bconfig, make(chan int))
	backend.Start()
	if err := db.Init("chatdata", backend.NodeID()); err != nil {
		log.Panicln(err)
	}
	config := server.HTTPConfig{Address: "127.0.0.1:0"}
	server.RunHTTPServer(backend, config)
}
