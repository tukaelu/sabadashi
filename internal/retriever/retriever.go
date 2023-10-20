package retriever

import (
	"context"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/tukaelu/sabadashi/internal/exporter"
)

func Retrieve(ctx *context.Context, metricName string, fn func() ([]mackerel.MetricValue, error)) (exporter.CsvRecords, error) {
	metricValues, err := fn()
	if err != nil {
		return nil, err
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
