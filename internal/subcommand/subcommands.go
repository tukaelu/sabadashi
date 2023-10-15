package subcommand

import (
	"github.com/urfave/cli/v2"

	"github.com/tukaelu/sabadashi/internal/subcommand/host"
)

func NewSubCommands() []*cli.Command {
	return []*cli.Command{
		host.NewHostSubcommand(),
	}
}
