package exporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStringArray(t *testing.T) {
	rec := &CsvRecord{
		Name:  "cpu.guest.percentage",
		Time:  1696113420, // 2023-10-01 07:37:00
		Value: 0,
	}

	normal := rec.ToStringArray(false)
	assert.Equal(t, "cpu.guest.percentage", normal[0], "Name of arrayed elements must match.")
	assert.Equal(t, "1696113420", normal[1], "Time of arrayed elements must match.")
	assert.Equal(t, "0", normal[2], "Value of arrayed elements must match.")

	friendly := rec.ToStringArray(true)
	assert.Equal(t, "2023-10-01 07:37:00 +0900 JST", friendly[0], "Datetime of arrayed elements must match.")
	assert.Equal(t, "cpu.guest.percentage", friendly[1], "Name of arrayed elements must match.")
	assert.Equal(t, "1696113420", friendly[2], "Time of arrayed elements must match.")
	assert.Equal(t, "0", friendly[3], "Value of arrayed elements must match.")
}
