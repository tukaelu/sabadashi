package subcommand

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/urfave/cli/v2"
)

type baseCommand struct {
	client       *mackerel.Client
	name         string
	from         int64
	to           int64
	granularity  string
	withFriendly bool
	rawFrom      string
	rawTo        string
}

func NewSubCommands() []*cli.Command {
	return []*cli.Command{
		NewHostSubcommand(),
		NewServiceSubCommand(),
	}
}
