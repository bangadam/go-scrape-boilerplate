package utils

import (
	"fmt"
	"time"
)

// string to float64
func StringToFloat64(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// string date to time.Time
func StringDateToTime(s string, layout string) time.Time {
	t, _ := time.Parse(layout, s)
	return t
}

// string to int64
func StringToInt64(s string) int64 {
	var i int64
	fmt.Sscanf(s, "%d", &i)
	return i
}
