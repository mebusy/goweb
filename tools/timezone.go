package tools

import (
    "time"
)

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
