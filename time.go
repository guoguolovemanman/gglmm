package gglmm

import (
	"time"
)

// ParseToTime 把值转换成时间
func ParseToTime(layout string, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

// ParseToTimes 把值转换成时间
func ParseToTimes(layout string, values []string) ([]time.Time, error) {
	result := make([]time.Time, 0)
	for _, value := range values {
		timeValue, err := ParseToTime(layout, value)
		if err != nil {
			return nil, err
		}
		result = append(result, timeValue)
	}
	return result, nil
}
