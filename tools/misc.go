package tools

import (
	"bytes"
	"encoding/json"
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

// will keep integer, rather than float64
func Struct2Json2Map( s interface{}  ) (map[string]interface{}, error ) {
    b, err := json.Marshal( s )
    if err != nil {
        return nil , err
    }
    m := map[string] interface{} {} 
    
    // simple json.UnMarshal to interface{} will generate float64 number
    // using json.NewDecoder.Decode instead
    d := json.NewDecoder( bytes.NewReader(b) )
    d.UseNumber() // keep integer, rather than float64
    err = d.Decode( &m )
    if err != nil {
        return nil , err
    }
    return m , nil
}



