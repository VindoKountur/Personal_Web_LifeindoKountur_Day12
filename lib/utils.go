package utils

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func StringInSlice(slc []string, str string) bool {
	for _, value := range slc {
		if value == str {
			return true
		}
	}
	return false
}

func CountDuration(start time.Time, end time.Time) string {
	timeDifferent := float64(end.Sub(start).Milliseconds())

	monthDistance := int(math.Floor(timeDifferent / (30 * 24 * 60 * 60 * 1000)))
	weekDistance := int(math.Floor(timeDifferent / (7 * 24 * 60 * 60 * 1000)))
	dayDistance := int(math.Floor(timeDifferent / (24 * 60 * 60 * 1000)))

	if monthDistance > 0 {
		str := strconv.Itoa(monthDistance) + " month"
		if monthDistance > 1 {
			return str + "s"
		}
		return str
	}
	if weekDistance > 0 {
		str := strconv.Itoa(weekDistance) + " week"
		if weekDistance > 1 {
			return str + "s"
		}
		return str
	}
	if dayDistance > 0 {
		str := strconv.Itoa(dayDistance) + " day"
		if dayDistance > 1 {
			return str + "s"
		}
		return str
	}
	return "cannot get duration"
}

func GetDateFormat(t time.Time) string {
	day := strconv.Itoa(t.Day())
	month := t.Month()
	year := strconv.Itoa(t.Year())
	result := fmt.Sprintf("%s %s %s", day, month, year)
	return result
}

func InputHtmlDateFormat(t time.Time) string {
	day := strconv.Itoa(t.Day())
	month := strconv.Itoa(int(t.Month()))
	year := strconv.Itoa(t.Year())
	result := fmt.Sprintf("%s-%02s-%02s", year, month, day)
	return result
}
