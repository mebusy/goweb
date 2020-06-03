package tools

import (
	"time"
)

func GetMillis() int64 {
	now := time.Now().UTC()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func GetSeconds() int64 {
	now := time.Now().UTC()
	return now.Unix()
}
