package main

import (
    "log"
    "github.com/mebusy/goweb/tools"
    "github.com/mebusy/goweb/encrypt"
)

func main() {
    log.Println( tools.UtctimestampToBeijingDate( tools.GetSeconds(), "2006-01-02 15:04:05" ) )
    log.Println( encrypt.SignSHA256withRSA( "a", "b", false ) )
    log.Println("done")
}
