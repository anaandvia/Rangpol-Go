package helper

import (
	"strings"
	"time"
)

func GetIndonesianDay(day time.Weekday) string {
	days := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}
	return days[day]
}

func toUpperCase(s interface{}) string {
	// Assert the type to string
	if str, ok := s.(string); ok {
		return strings.ToUpper(str)
	}
	return ""
}
