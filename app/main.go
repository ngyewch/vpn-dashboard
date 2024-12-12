package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	version         string
	commit          string
	commitTimestamp string

	flagListenAddr = &cli.StringFlag{
		Name:    "listen-addr",
		Usage:   "listen addr",
		Value:   ":8080",
		EnvVars: []string{"LISTEN_ADDR"},
	}

	app = &cli.App{
		Name:    "vpn-dashboard",
		Usage:   "VPN dashboard",
		Version: version,
		Action:  nil,
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
	cli.VersionPrinter = func(cCtx *cli.Context) {
		var parts []string
		if version != "" {
			parts = append(parts, fmt.Sprintf("version=%s", version))
		}
		if commit != "" {
			parts = append(parts, fmt.Sprintf("commit=%s", commit))
		}
		if commitTimestamp != "" {
			formattedCommitTimestamp := func(commitTimestamp string) string {
				epochSeconds, err := strconv.ParseInt(commitTimestamp, 10, 64)
				if err != nil {
					return ""
				}
				t := time.Unix(epochSeconds, 0)
				return t.Format(time.RFC3339)
			}(commitTimestamp)
			if formattedCommitTimestamp != "" {
				parts = append(parts, fmt.Sprintf("commitTimestamp=%s", formattedCommitTimestamp))
			}
		}
		fmt.Println(strings.Join(parts, " "))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
