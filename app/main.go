package main

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

var (
	version string

	flagListenAddr = &cli.StringFlag{
		Name:    "listen-addr",
		Usage:   "listen addr",
		Value:   ":8080",
		Sources: cli.EnvVars("LISTEN_ADDR"),
	}

	app = &cli.Command{
		Name:    "vpn-dashboard",
		Usage:   "VPN dashboard",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "serve",
				Action: doServe,
				Flags: []cli.Flag{
					flagListenAddr,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
