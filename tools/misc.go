package tools

import (
    "reflect"
    "strings"
)


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

func Lower1stCharOfStringArray ( arr []string ) []string {
    arr_lower := []string {}
    for _, str := range arr {
        if str == "" {
            arr_lower = append( arr_lower, str  )
        } else {
            arr_lower = append( arr_lower, strings.ToLower(str[0:1])+ str[1:]  )
        }
    }
    return arr_lower
}

func GetStructFieldNum( u interface{} ) int {
    t := reflect.TypeOf(u)
    if t.Kind() == reflect.Ptr { // if it is pointer
        t = t.Elem() // get the actually type
    }
    return t.NumField()
}



