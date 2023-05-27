package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
			Name:  "config",
			Value: "config-p2p.ini",
		},
	},
	Action: p2pchat,
}

func p2pchat(ctx *cli.Context) error {
	cfgfile := ctx.String("config")
	if _, err := os.Stat(cfgfile); err != nil {
		return errors.Join(
			fmt.Errorf("cannot load config file: %s", cfgfile),
			err,
		)
	}
	config := config.LoadINIConfig(cfgfile)
	backend := api.NewBackend(config.Backend)
	backend.Start()
	defer func() {
		backend.Stop()
		log.Println("backend closed")
	}()
	if err := db.Init(config.DataDir, backend.NodeID()); err != nil {
		log.Panicln(err)
	}
	httpserver, _, err := server.StartHTTPServer(backend, config.Http)
	if err != nil {
		return err
	}
	defer func() {
		httpserver.Shutdown(context.Background())
		log.Println("http server closed")
	}()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	return nil
}

func main() {
	SetLogLevel(LvlFromString("debug"))
	if err := app.Run(os.Args); err != nil {
		log.Panicln(err)
	}
}
