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

type serviceCommand struct {
	baseCommand
	name         string
	withExternal bool
}

func (c *serviceCommand) run(ctx *cli.Context) error {

	exportDir, err := fileutil.CreateExportDir(c.name, c.rawFrom, c.rawTo)
	if err != nil {
		return err
	}

	interval := converter.GranularityToInterval(c.granularity)
	attempts := (c.to-c.from)/interval + 1

	metricNames, err := c.client.ListServiceMetricNames(c.name)
	if err != nil {
		return err
	}

	if c.withExternal {
		externalMetricNames, err := c.listExternalMonitorMetricNames()
		if err != nil {
			return err
		}
		metricNames = append(metricNames, externalMetricNames...)
	}

	if len(metricNames) == 0 {
		fileutil.RemoveDir(exportDir)
		return fmt.Errorf("There was no service metric available for export.")
	}

	progress := progressbar.Default(attempts * int64(len(metricNames)))

	ch := make(chan exporter.Result, definition.CONCURRENT_FILE_OPERATION)
	defer close(ch)

	go func(ch chan exporter.Result) {
		for {
			select {
			case res := <-ch:
				if err := fileutil.WriteFile(exportDir, res.MetricName, "csv", res.Records, c.withFriendly); err != nil {
					fmt.Printf("Failed to write CSV file. (reason: %s)\n", err)
					return
				}
			case <-ctx.Done():
				fmt.Println("Operation canceled.")
				os.Exit(125)
			}
		}
	}(ch)

	from := c.from
	to := int64(0)
	wg := &sync.WaitGroup{}
	for i := int64(0); i < attempts; i++ {
		to = from + interval
		if to > c.to {
			to = c.to
		}

		for _, metricName := range metricNames {
			wg.Add(1)
			go func(ctx *context.Context, metricName string) {
				defer wg.Done()

				metrics, err := retriever.Retrieve(ctx, metricName, func() ([]mackerel.MetricValue, error) {
					metricValues, err := c.client.FetchServiceMetricValues(c.name, metricName, from, to)
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

func (c *serviceCommand) listExternalMonitorMetricNames() ([]string, error) {
	monitors, err := c.client.FindMonitors()
	if err != nil {
		return nil, err
	}

	var metricNames []string
	for _, monitor := range monitors {
		if monitor.MonitorType() != "external" {
			continue
		}
		em := monitor.(*mackerel.MonitorExternalHTTP)
		if em.Service != c.name {
			continue
		}
		metricNames = append(metricNames, fmt.Sprintf("__externalhttp.responsetime.%s", monitor.MonitorID()))
	}
	return metricNames, nil
}
