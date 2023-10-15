package converter

import (
	"time"

	"github.com/tukaelu/sabadashi/internal/definition"
)

// ToStartDayOfUnixtime returns the start date and time of the specified day.
func ToStartDayOfUnixtime(ds string) int64 {
	p, _ := time.Parse(definition.DATEFORMAT, ds)
	return time.Date(p.Year(), p.Month(), p.Day(), 0, 0, 0, 0, time.Local).Unix()
}

// ToEndDayOfUnixtime returns the last date and time of the specified day.
func ToEndDayOfUnixtime(ds string) int64 {
	p, _ := time.Parse(definition.DATEFORMAT, ds)
	return time.Date(p.Year(), p.Month(), p.Day(), 23, 59, 59, 0, time.Local).Unix()
}

// GranularityToInterval converts from granularity to Unixtime acquisition interval.
func GranularityToInterval(s string) int64 {
	switch s {
	case "1m":
		return definition.GRANULARITY_1M_INTERVAL
	case "5m":
		return definition.GRANULARITY_5M_INTERVAL
	case "10m":
		return definition.GRANULARITY_10M_INTERVAL
	case "1h":
		return definition.GRANULARITY_1H_INTERVAL
	case "2h":
		return definition.GRANULARITY_2H_INTERVAL
	case "4h":
		return definition.GRANULARITY_4H_INTERVAL
	case "1d":
		return definition.GRANULARITY_1D_INTERVAL
	default:
		return definition.GRANULARITY_1M_INTERVAL
	}
}
