package tools

import (
	"time"
    "reflect"
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

var secondsEastOfUTC = int((8 * time.Hour).Seconds())
var beijing = time.FixedZone("Beijing Time", secondsEastOfUTC)

// datetime in `format ` , to UTC timestamp
func UtctimestampFromBeijingDate( date string, format string ) (int64,error) {
    // RFC3339     := "2006010215"
    time_target, err  := time.ParseInLocation( format,   date , beijing )
    if err != nil {
        return  -1, err
    }
    return time_target.UTC().Unix(), nil
}

func UtctimestampToBeijingDate( seconds int64, format string ) string {
    if seconds == 0 {
        return "N/A"
    }
    // FORMAT := "2006-01-02 15:04"

    t := time.Unix( seconds,0 ).UTC()
    return t.In(beijing).Format( format )
}

func GetStructFields( u interface{} ) []string {
    fields := []string{}
    t := reflect.TypeOf(u)
    if t.Kind() == reflect.Ptr { // if it is pointer
        t = t.Elem() // get the actually type
    }
    for i, n := 0, t.NumField(); i < n; i++ {
        f := t.Field(i)
        fields = append(fields, f.Name)
    }
    return fields
}

