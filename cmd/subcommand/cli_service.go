package subcommand

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/urfave/cli/v2"

	"github.com/tukaelu/sabadashi/internal/converter"
	"github.com/tukaelu/sabadashi/internal/validator"
)

func NewServiceSubCommand() *cli.Command {
	return &cli.Command{
		Name:  "service",
		Usage: "Retrieves service metrics",
		Action: func(ctx *cli.Context) error {

			from := converter.ToStartDayOfUnixtime(ctx.String("from"))
			to := converter.ToEndDayOfUnixtime(ctx.String("to"))

			if err := validator.ValidateTheDateBeforeAndAfter(from, to); err != nil {
				return err
			}
			if err := validator.ValidateMetricRetantionPeriod(from, to); err != nil {
				return err
			}

			cmd := &serviceCommand{
				baseCommand: baseCommand{
					client:       mackerel.NewClient(ctx.String("apikey")),
					from:         from,
					to:           to,
					granularity:  ctx.String("granularity"),
					withFriendly: ctx.Bool("with-friendly-date-format"),
					rawFrom:      ctx.String("from"),
					rawTo:        ctx.String("to"),
				},
				name:         ctx.String("name"),
				withExternal: ctx.Bool("with-external-monitors"),
			}

			return cmd.run(ctx)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the service from which to retrieve metric",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "from",
				Usage:    "Specify the date to start retrieving metrics in YYYYMMDD format. (e.g. 20230101)",
				Required: true,
				Action: func(ctx *cli.Context, s string) error {
					return validator.ValidateDateFormat(s)
				},
			},
			&cli.StringFlag{
				Name:     "to",
				Usage:    "Specify the date to end retrieving metrics in YYYYMMDD format. (e.g. 20231231)",
				Required: true,
				Action: func(ctx *cli.Context, s string) error {
					return validator.ValidateDateFormat(s)
				},
			},
			&cli.StringFlag{
				Name:        "granularity",
				Aliases:     []string{"g"},
				Usage:       "Specify the granularity of metric data. Choose from 1m, 5m, 10m, 1h, 2h, 4h or 1d.",
				Value:       "1m",
				DefaultText: "1m",
				Action: func(ctx *cli.Context, s string) error {
					return validator.ValidateGranularity(s)
				},
			},
			&cli.BoolFlag{
				Name:    "with-friendly-date-format",
				Aliases: []string{"f"},
				Usage:   "If this flag is enabled, an additional column with a friendly date format is output at the beginning of the CSV line.",
			},
			&cli.BoolFlag{
				Name:    "with-external-monitors",
				Aliases: []string{"e"},
				Usage:   "If this flag is enabled, it also includes the metric measured in the external monitoring.",
			},
		},
	}
}
