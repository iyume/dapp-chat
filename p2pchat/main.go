package main

import (
	"log"

	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/config"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/server"
)

func main() {
	SetLogLevel(LvlFromString("debug"))

	config := config.LoadINIConfig("config-p2p.ini")
	backend := api.NewBackend(config.Backend, make(chan int))
	backend.Start()
	if err := db.Init("chatdata", backend.NodeID()); err != nil {
		log.Panicln(err)
	}
	server.RunHTTPServer(backend, config.Http)
}
