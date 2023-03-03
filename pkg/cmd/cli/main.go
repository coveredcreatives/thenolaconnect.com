package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start running server instance",
				Action: func(ctx *cli.Context) error {
					return Server(ctx)
				},
			},
			{
				Name:  "db",
				Usage: "start empty database instance",
				Action: func(ctx *cli.Context) error {
					terminate, err := ContainerDatabase(ctx)
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
					return SynchronizeDB(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
