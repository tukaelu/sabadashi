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

func TestToStringt(t *testing.T) {
	namePtn := "cpu.idle.percentage"
	timePtn := int64(1696086000)
	valuePtn1 := interface{}(1316306944)
	valuePtn2 := interface{}(23569.216666666667)
	valuePtn3 := interface{}(5.0999080344452805)
	assert.Equal(t, namePtn, toString(namePtn), "convert namePtn to string")
	assert.Equal(t, "1696086000", toString(timePtn), "convert timePtn to string")
	assert.Equal(t, "1316306944", toString(valuePtn1), "convert valuePtn1 to string (Not to be in exponential notation)")
	assert.Equal(t, "23569.216666666667", toString(valuePtn2), "convert valuePtn2 to string (Not to be in exponential notation)")
	assert.Equal(t, "5.0999080344452805", toString(valuePtn3), "convert valuePtn3 to string (Not to be in exponential notation)")
}
