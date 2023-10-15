package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tukaelu/sabadashi/internal/definition"
)

func TestToStartDayOfUnixtime(t *testing.T) {
	expected := int64(1696086000) // 2023-10-01 00:00:00
	assert.Equal(t, expected, ToStartDayOfUnixtime("20231001"), "nil returned for correct date")
}

func TestToEndDayOfUnixtime(t *testing.T) {
	expected := int64(1696172399) // 2023-10-01 23:59:59
	assert.Equal(t, expected, ToEndDayOfUnixtime("20231001"), "nil returned for correct date")
}

func TestGranularityToInterval(t *testing.T) {
	assert.Equal(t, int64(definition.GRANULARITY_1M_INTERVAL), GranularityToInterval("1m"), "The acquisition period must correspond to 1-min granularity.")
	assert.Equal(t, int64(definition.GRANULARITY_5M_INTERVAL), GranularityToInterval("5m"), "The acquisition period must correspond to 5-min granularity.")
	assert.Equal(t, int64(definition.GRANULARITY_10M_INTERVAL), GranularityToInterval("10m"), "The acquisition period must correspond to 10-min granularity.")
	assert.Equal(t, int64(definition.GRANULARITY_1H_INTERVAL), GranularityToInterval("1h"), "The acquisition period must correspond to 1-hour granularity.")
	assert.Equal(t, int64(definition.GRANULARITY_2H_INTERVAL), GranularityToInterval("2h"), "The acquisition period must correspond to 2-hour granularity.")
	assert.Equal(t, int64(definition.GRANULARITY_4H_INTERVAL), GranularityToInterval("4h"), "The acquisition period must correspond to 4-hour granularity.")
	assert.Equal(t, int64(definition.GRANULARITY_1D_INTERVAL), GranularityToInterval("1d"), "The acquisition period must correspond to 1-day granularity.")
	assert.Equal(t, int64(definition.GRANULARITY_1M_INTERVAL), GranularityToInterval("2d"), "In the case of unexpected parameter, the acquisition period should correspond to 1-min granularity.")
}
