package gglmm

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	timeString := "2006-01-02T15:04:05+08:00"
	timeValue, err := ParseToTime(time.RFC3339, timeString)
	if err != nil {
		t.Fail()
	}
	t.Log(timeValue)
}

func TestParseTimes(t *testing.T) {
	timeStrings := []string{"2006-01-02T15:04:05+08:00", "2006-01-02T15:04:05+09:00"}
	timeValues, err := ParseToTimes(time.RFC3339, timeStrings)
	if err != nil {
		t.Fail()
	}
	t.Log(timeValues)
}
