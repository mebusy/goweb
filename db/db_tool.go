package db

import (
    "github.com/mebusy/goweb/tools"
)

func GetDBFieldListFromStrunct( u interface{}, bLower1st bool ) []string {
    if bLower1st {
        return tools.Lower1stCharOfStringArray( tools.GetStructFields( u ) )
    } else {
        return tools.GetStructFields( u )
    }
}
