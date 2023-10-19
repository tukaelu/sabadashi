package subcommand

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/urfave/cli/v2"

	"github.com/tukaelu/sabadashi/internal/converter"
	"github.com/tukaelu/sabadashi/internal/validator"
)

func NewHostSubcommand() *cli.Command {
	return &cli.Command{
		Name:  "host",
		Usage: "Retrieves host metrics",
		Action: func(ctx *cli.Context) error {

			// [TODO] Check duplicate processes on the same host

			from := converter.ToStartDayOfUnixtime(ctx.String("from"))
			to := converter.ToEndDayOfUnixtime(ctx.String("to"))

			if err := validator.ValidateTheDateBeforeAndAfter(from, to); err != nil {
				return err
			}
			if err := validator.ValidateMetricRetantionPeriod(from, to); err != nil {
				return err
			}

			cmd := &hostCommand{
				baseCommand: baseCommand{
					client:       mackerel.NewClient(ctx.String("apikey")),
					from:         from,
					to:           to,
					granularity:  ctx.String("granularity"),
					withFriendly: ctx.Bool("with-friendly-date-format"),
					rawFrom:      ctx.String("from"),
					rawTo:        ctx.String("to"),
				},
				host: ctx.String("host"),
			}
			return doHost(ctx, cmd)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"H"},
				Usage:    "ID of the host from which to retrieve metric",
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
		},
	}
}
