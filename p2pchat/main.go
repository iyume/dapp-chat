package main

import (
	"log"
	"os"

	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/config"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/server"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "p2pchat",
	Usage: "go run ./p2pchat --config FILE [...]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Required: true,
		},
		&cli.StringFlag{
			Name: "data",
		},
	},
	Action: func(ctx *cli.Context) error {
		cfgfile := ctx.String("config")
		if _, err := os.Stat(cfgfile); err != nil {
			log.Panicln(err)
		}
		datadir := ctx.String("data")
		if datadir != "" {
			if _, err := os.Stat(datadir); err != nil {
				log.Panicln(err)
			}
		} else {
			datadir = "chatdata"
		}
		config := config.LoadINIConfig(cfgfile)
		backend := api.NewBackend(config.Backend, make(chan int))
		backend.Start()
		if err := db.Init(datadir, backend.NodeID()); err != nil {
			log.Panicln(err)
		}
		server.RunHTTPServer(backend, config.Http)
		return nil
	},
}

func main() {
	SetLogLevel(LvlFromString("debug"))
	if err := app.Run(os.Args); err != nil {
		log.Panicln(err)
	}
}
