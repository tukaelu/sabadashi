package exporter

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

type CsvRecord struct {
	Name  string
	Time  int64
	Value interface{}
}

type CsvRecords []*CsvRecord

type Result struct {
	MetricName string
	Records    CsvRecords
}

// ToStringArray converts a structure to an array and returns it.
// If the option is enabled, it adds a column at the beginning of the line with a date that is friendly format.
func (c *CsvRecord) ToStringArray(friendly bool) []string {
	if friendly {
		return []string{
			// [TODO] considering the load of conversion, it might be better to use cache or make it optional.
			time.Unix(c.Time, 0).Local().String(),
			c.Name,
			strconv.FormatInt(c.Time, 10),
			toString(c.Value),
		}
	} else {
		return []string{
			c.Name,
			strconv.FormatInt(c.Time, 10),
			toString(c.Value),
		}
	}
}

func toString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case int, int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case []byte:
		return *(*string)(unsafe.Pointer(&v))
	default:
		return ""
	}
}
