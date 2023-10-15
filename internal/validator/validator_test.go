package validator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDateFormat(t *testing.T) {
	assert.Equal(t, nil, ValidateDateFormat("20231001"), "nil returned for correct date")
	assert.NotEqual(t, nil, ValidateDateFormat("2023-10-01"), "Not nil because it is not in the expected date format")
}

func TestValidateMetricRetantionPeriod(t *testing.T) {
	from := int64(1672498800)       // 2023-01-01 00:00:00
	to := int64(1712242799)         // 2024-04-04 23:59:59
	toOver1day := int64(1712329199) // 2024-04-05 23:59:59
	assert.Equal(t, nil, ValidateMetricRetantionPeriod(from, to), "within the retention period of the metric.")
	assert.Equal(
		t,
		fmt.Errorf("specified period exceeds the data retention period of 460 days. (Your specification is 461 days.)"),
		ValidateMetricRetantionPeriod(from, toOver1day),
		"specified for a metric retention period of one day beyond that specified.",
	)
}

func TestValidateGranularity(t *testing.T) {
	expectedErr := fmt.Errorf("unsupported granularity pattern 1w has been set")
	assert.Equal(t, nil, ValidateGranularity("1m"), "The specified granularity that is supported.")
	assert.Equal(t, nil, ValidateGranularity("5m"), "The specified granularity that is supported.")
	assert.Equal(t, nil, ValidateGranularity("10m"), "The specified granularity that is supported.")
	assert.Equal(t, nil, ValidateGranularity("1h"), "The specified granularity that is supported.")
	assert.Equal(t, nil, ValidateGranularity("2h"), "The specified granularity that is supported.")
	assert.Equal(t, nil, ValidateGranularity("4h"), "The specified granularity that is supported.")
	assert.Equal(t, nil, ValidateGranularity("1d"), "The specified granularity that is supported.")
	assert.Equal(t, expectedErr, ValidateGranularity("1w"), "Unsupported granularity was specified.")
}