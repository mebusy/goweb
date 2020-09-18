package tools

import (
    "time"
)

func GetMillis() int64 {
    now := time.Now().UTC()
    nanos := now.UnixNano()
    millis := nanos / 1000000
    return millis + ( secondsDelta * 1000 )
}

func GetSeconds() int64 {
    now := time.Now().UTC()
    return now.Unix() + secondsDelta
}

var secondsDelta int64 = 0

func DEBUG_resetTime() {
    secondsDelta = 0
}

func DEBUG_setDelta( delta int64) {
    secondsDelta = delta
}

func DEBUG_GetDelta( ) int64  {
    return secondsDelta
}



