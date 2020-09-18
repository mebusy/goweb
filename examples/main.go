package main

import (
    "log"
    "github.com/mebusy/goweb/tools"
)

func main() {
    log.Println( tools.UtctimestampToBeijingDate( tools.GetSeconds(), "2006-01-02 15:04:05" ) )
    log.Println("done")
}
