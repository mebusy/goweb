package main

import (
    "log"
    "github.com/mebusy/goweb/tools"
)

func main() {
    log.Println( tools.UtctimestampToBeijingDate( 1, "2006-01-02" ) )
    log.Println("done")
}
