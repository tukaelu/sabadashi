package host

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"

	"github.com/tukaelu/sabadashi/internal/converter"
	"github.com/tukaelu/sabadashi/internal/definition"
	"github.com/tukaelu/sabadashi/internal/exporter"
)

type hostCommand struct {
	client      *mackerel.Client
	host        string
	from        int64
	to          int64
	granularity string
	friendly    bool
	rawFrom     string
	rawTo       string
}

func run(ctx *cli.Context, h *hostCommand) error {

	// [FIXME] If it is run on the same host and time period, we will tentatively treat it as an error, so please back it up and run it again.
	creationDir := getCreationDirPath(h.host, h.rawFrom, h.rawTo)
	if _, err := os.Stat(creationDir); os.IsNotExist(err) {
		if err := os.MkdirAll(creationDir, 0755); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("export directory '%s' is already exists, so please back it up and run it again", creationDir)
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
				if err := writeCsvFile(creationDir, res.MetricName, res.Records, h.friendly); err != nil {
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

				metrics, err := retrieve(ctx, h.client, h.host, metricName, from, to)
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

func retrieve(ctx *context.Context, client *mackerel.Client, host, metricName string, from, to int64) (exporter.CsvRecords, error) {
	// [TODO] if err do retry
	metricValues, err := client.FetchHostMetricValues(host, metricName, from, to)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to retrieve metrics (metric: %s, from: %d, to: %d) (reason: %s)", metricName, from, to, err)
	}

	var records exporter.CsvRecords
	for _, metricValue := range metricValues {
		record := &exporter.CsvRecord{
			Name:  metricName,
			Time:  metricValue.Time,
			Value: metricValue.Value,
		}
		records = append(records, record)
	}

	return records, nil
}

func writeCsvFile(creationDir, metricName string, metrics exporter.CsvRecords, friendly bool) error {
	outfile := filepath.Join(
		creationDir,
		fmt.Sprintf("%s.csv", metricName),
	)

	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// [TODO] The standard Go encoding/csv library may not be able to handle CSV escape, therefore, it may be a good idea to introduce a library for this purpose.
	cw := csv.NewWriter(f)
	defer cw.Flush()

	for _, metric := range metrics {
		if err := cw.Write(metric.ToStringArray(friendly)); err != nil {
			fmt.Printf("write csv file failed (reason: %s)", err)
		}
	}

	return nil
}

func getCreationDirPath(hostId, from, to string) string {
	cwd, _ := os.Getwd()
	return filepath.Join(
		cwd,
		hostId,
		fmt.Sprintf("%s_%s", from, to),
	)
}
