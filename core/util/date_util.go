package util

import (
	"time"
)

const (
	YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05"
	YYYY_MM_DD          = "2006-01-02"
	YYYY_MM             = "2006-01"
	YYYYMMDDHHMMSS      = "20060102150405"
	YYYYMMDD_HH         = "20060102_15"
)

func GetTime() time.Time {
	time, _ := time.ParseInLocation(YYYYMMDD_HH, time.Now().Format(YYYYMMDD_HH), time.Local)
	return time
}

func GetCurrentDateStr(pattern string) string {
	return time.Now().Format(pattern)
}

func GetDateStr(time time.Time, pattern string) string {
	return time.Format(pattern)
}

func Str2Date(dateStr, pattern string) time.Time {
	time, _ := time.ParseInLocation(pattern, dateStr, time.Local)
	return time
}

func AddMonth(time time.Time, i int) time.Time {
	return time.AddDate(0, i, 0)
}

func GetMonthActualMaximum(time time.Time) int {

	var testRunNian = func() int {
		year, _, _ := time.Date()
		if year%4 == 0 && year%100 != 0 {
			return 29
		} else {
			return 28
		}
	}

	switch time.Month() {
	case 1:
		return 31
	case 2:
		return testRunNian()
	case 3:
		return 31
	case 4:
		return 30
	case 5:
		return 31
	case 6:
		return 30
	case 7:
		return 31
	case 8:
		return 31
	case 9:
		return 30
	case 10:
		return 31
	case 11:
		return 30
	case 12:
		return 31
	default:
		panic("error")
	}
}
