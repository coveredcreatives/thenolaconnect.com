package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/devtools"

	cli "github.com/urfave/cli/v2"
)

func main() {
	v, err := devtools.Config()
	if err != nil {
		log.WithError(err).Error("failed to load config")
		os.Exit(1)
	}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start running server instance",
				Action: func(ctx *cli.Context) error {
					return Server(ctx, v)
				},
			},
			{
				Name:  "db",
				Usage: "start empty database instance",
				Action: func(ctx *cli.Context) error {
					terminate, err := ContainerDatabase(ctx, v)
					if err != nil {
						return err
					}
					defer terminate()
					quitChannel := make(chan os.Signal, 1)
					signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
					<-quitChannel
					return nil
				},
			},
			{
				Name:  "sync",
				Usage: "sync database with latest status of orders form response",
				Action: func(ctx *cli.Context) error {
					return SynchronizeDB(ctx, v)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error("exiting")
		os.Exit(1)
	}

}
