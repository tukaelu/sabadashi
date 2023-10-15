package validator

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"

	"github.com/tukaelu/sabadashi/internal/definition"
)

// ValidateDateFormat ...
func ValidateDateFormat(ds string) error {
	_, err := time.Parse(definition.DATEFORMAT, ds)
	if err != nil {
		return err
	}
	return nil
}

// ValidateTheDateBeforeAndAfter ...
func ValidateTheDateBeforeAndAfter(from, to int64) error {
	if from > to {
		return fmt.Errorf("the start and end dates are reversed")
	}
	return nil
}

// ValidateMetricRetantionPeriod ...
// If the specified period exceeds the retention period limit of the Standard Plan, an error will occure.
// see. https://en.mackerel.io/pricing
func ValidateMetricRetantionPeriod(from, to int64) error {
	period := to - from
	if period > definition.MAX_DATA_RETENTION_INTERVAL {
		return fmt.Errorf(
			"specified period exceeds the data retention period of %d days. (Your specification is %d days.)",
			definition.MAX_DATA_RETENTION_DAYS,
			(period+1)/(60*60*24),
		)
	}
	return nil
}

// ValidateGranularity ...
func ValidateGranularity(s string) error {
	valid := []string{"1m", "5m", "10m", "1h", "2h", "4h", "1d"}
	if !slices.Contains(valid, s) {
		return fmt.Errorf("unsupported granularity pattern %s has been set", s)
	}
	return nil
}
