package subcommand

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"

	"github.com/tukaelu/sabadashi/internal/converter"
	"github.com/tukaelu/sabadashi/internal/definition"
	"github.com/tukaelu/sabadashi/internal/exporter"
	"github.com/tukaelu/sabadashi/internal/fileutil"
	"github.com/tukaelu/sabadashi/internal/retriever"
)

type hostCommand struct {
	baseCommand
	host string
}

func doHost(ctx *cli.Context, h *hostCommand) error {

	exportDir, err := fileutil.CreateExportDir(h.host, h.rawFrom, h.rawTo)
	if err != nil {
		return err
	}

	interval := converter.GranularityToInterval(h.granularity)
	attempts := (h.to-h.from)/interval + 1

	metricNames, err := h.client.ListHostMetricNames(h.host)
	if err != nil {
		return err
	}

	progress := progressbar.Default(attempts * int64(len(metricNames)))

	ch := make(chan exporter.Result, definition.CONCURRENT_FILE_OPERATION)
	defer close(ch)

	go func(ch chan exporter.Result) {
		for {
			select {
			case res := <-ch:
				if err := fileutil.WriteFile(exportDir, res.MetricName, "csv", res.Records, h.withFriendly); err != nil {
					fmt.Printf("Failed to write CSV file. (reason: %s)\n", err)
					return
				}
			case <-ctx.Done():
				fmt.Println("Operation canceled.")
				os.Exit(125)
			}
		}
	}(ch)

	from := h.from
	to := int64(0)
	wg := &sync.WaitGroup{}
	for i := int64(0); i < attempts; i++ {
		to = from + interval
		if to > h.to {
			to = h.to
		}

		for _, metricName := range metricNames {
			wg.Add(1)
			go func(ctx *context.Context, metricName string) {
				defer wg.Done()

				metrics, err := retriever.Retrieve(ctx, metricName, func() ([]mackerel.MetricValue, error) {
					metricValues, err := h.client.FetchHostMetricValues(h.host, metricName, h.from, h.to)
					if err != nil {
						return nil, fmt.Errorf("[ERROR] failed to retrieve metrics (metric: %s, from: %d, to: %d) (reason: %s)", metricName, from, to, err)
					}
					return metricValues, nil
				})

				if err != nil {
					// skip retrieve metrics...
					fmt.Sprintln(err)
				}

				ch <- exporter.Result{
					MetricName: metricName,
					Records:    metrics,
				}
			}(&ctx.Context, metricName)

			time.Sleep(350 * time.Millisecond)

			_ = progress.Add(1)
		}
		wg.Wait()

		from += interval
		_ = progress.Add(1)
	}

	return nil
}
