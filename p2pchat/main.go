package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/server"
)

func main() {
	SetLogLevel(LvlFromString("debug"))

	bconfig := api.DefaultBackendConfig
	bconfig.Key, _ = crypto.GenerateKey()
	backend := api.NewBackend(bconfig, make(chan int))
	backend.Start()
	config := server.HTTPConfig{Address: ":0"}
	server.RunHTTPServer(backend, config)
}
