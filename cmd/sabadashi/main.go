package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/tukaelu/sabadashi/internal/subcommand"
)

var version string
var revision string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	defer cancel()

	if err := NewCliApp().RunContext(ctx, os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func NewCliApp() *cli.App {
	return &cli.App{
		Name:    "sabadashi",
		Usage:   "Retrieves metrics stored in Mackerel",
		Version: fmt.Sprintf("%s (rev.%s)", version, revision),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "apikey",
				Usage:    "API key to access the organization",
				EnvVars:  []string{"MACKEREL_APIKEY"},
				Required: true,
			},
		},
		Commands: subcommand.NewSubCommands(),
	}
}
