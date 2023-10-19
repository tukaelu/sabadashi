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

func doService(ctx *cli.Context, s *serviceCommand) error {

	exportDir, err := fileutil.CreateExportDir(s.name, s.rawFrom, s.rawTo)
	if err != nil {
		return err
	}

	interval := converter.GranularityToInterval(s.granularity)
	attempts := (s.to-s.from)/interval + 1

	metricNames, err := s.client.ListServiceMetricNames(s.name)
	if err != nil {
		return err
	}

	if s.withExternal {
		externalMetricNames, err := listExternalMonitorMetricNames(s)
		if err != nil {
			return err
		}
		metricNames = append(metricNames, externalMetricNames...)
	}

	if len(metricNames) == 0 {
		return fmt.Errorf("There was no service metric available for export.")
	}

	progress := progressbar.Default(attempts * int64(len(metricNames)))

	ch := make(chan exporter.Result, definition.CONCURRENT_FILE_OPERATION)
	defer close(ch)

	go func(ch chan exporter.Result) {
		for {
			select {
			case res := <-ch:
				if err := fileutil.WriteFile(exportDir, res.MetricName, "csv", res.Records, s.withFriendly); err != nil {
					fmt.Printf("Failed to write CSV file. (reason: %s)\n", err)
					return
				}
			case <-ctx.Done():
				fmt.Println("Operation canceled.")
				os.Exit(125)
			}
		}
	}(ch)

	from := s.from
	to := int64(0)
	wg := &sync.WaitGroup{}
	for i := int64(0); i < attempts; i++ {
		to = from + interval
		if to > s.to {
			to = s.to
		}

		for _, metricName := range metricNames {
			wg.Add(1)
			go func(ctx *context.Context, metricName string) {
				defer wg.Done()

				metrics, err := retriever.Retrieve(ctx, metricName, func() ([]mackerel.MetricValue, error) {
					metricValues, err := s.client.FetchServiceMetricValues(s.name, metricName, s.from, s.to)
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

func listExternalMonitorMetricNames(s *serviceCommand) ([]string, error) {
	monitors, err := s.client.FindMonitors()
	if err != nil {
		return nil, err
	}

	var metricNames []string
	for _, monitor := range monitors {
		if monitor.MonitorType() != "external" {
			continue
		}
		em := monitor.(*mackerel.MonitorExternalHTTP)
		if em.Service != s.name {
			continue
		}
		metricNames = append(metricNames, fmt.Sprintf("__externalhttp.responsetime.%s", monitor.MonitorID()))
	}
	return metricNames, nil
}
