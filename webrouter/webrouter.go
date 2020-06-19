package webrouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
)

var domain = "localhost:3000"

type API struct {
	Method  string
	Path    string
	Func    func(http.ResponseWriter, *http.Request)
	Params  interface{}
	Example string
	Desc    string
}

/*
Params example:
type ParamShark struct {
    Shark string
    Equipment string  `param:"optional"`
}
//*/

var api_doc = `
server response data :
    server responds a json string, like: {"data":{...},"err":"...","errcode":xxx} 
        it always has a 'data' filed which contains important data, and a 'errcode' field (-1 means no error) 
        the 'err' filed appears only if an error occur. and sometimes along with an extra 'errinfo' field for more details.
    client should always check whether 'err' field exists,  
        if the 'err' field do exist, client should do error handle.

API infomations: 
`

var api_struct_map = make(map[string]interface{})

func KeyFromRequestURI(uri string) string {
	idx := strings.Index(uri[1:], "/")
	if idx == -1 {
		return uri
	}
	return uri[:idx+1]
}

func mappingApiStruct(uri string, v interface{}, method string) {
	key := KeyFromRequestURI(uri)
	if key != "" {
		if method == "POST" {
			api_struct_map[key] = v
		}
		// fmt.Println( key )
	}
}

var (
	ReflectTypeInt    = reflect.TypeOf(int(1)) // TODO
	ReflectTypeString = reflect.TypeOf("")
)

// 将 struct 中的 int 字段，设置为 math.MinInt32
func initIntField(u interface{}) {
	t := reflect.TypeOf(u)
	if t.Kind() == reflect.Ptr { // 是指针
		t = t.Elem() // 进而获取 目标类型
	}
	v := reflect.ValueOf(u)
	if v.Kind() == reflect.Ptr { // 是指针
		v = v.Elem() // 进而获取 目标类型
	}

	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		if f.Type == ReflectTypeInt {
			v.FieldByName(f.Name).SetInt(math.MinInt32)
		}
	}
}

// 遍历 struct中的字段，如果是 初始值, 则报错
func verifyParams(u interface{}) error {
	t := reflect.TypeOf(u)
	if t.Kind() == reflect.Ptr { // 是指针
		t = t.Elem() // 进而获取 目标类型
	}
	v := reflect.ValueOf(u)
	if v.Kind() == reflect.Ptr { // 是指针
		v = v.Elem() // 进而获取 目标类型
	}
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		if tag := f.Tag.Get("param"); tag == "optional" {
			continue
		}
		switch f.Type {
		case ReflectTypeInt:
			// log.Println( f.Name ,  v.FieldByName( f.Name ).Int()   )
			if v.FieldByName(f.Name).Int() == math.MinInt32 {
				return errors.New("invalid field: " + f.Name)
			}
		case ReflectTypeString:
			if v.FieldByName(f.Name).String() == "" {
				return errors.New("invalid field: " + f.Name)
			}
		}
	}

	return nil
}

var funcErrResponse func(http.ResponseWriter, error)

func SetErrResponseFunc(f func(http.ResponseWriter, error)) {
	funcErrResponse = f
}

// type MiddlewareFunc func(http.Handler) http.Handler
// must register before  body-reading method
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		data, err := httputil.DumpRequest(r, true)

		if err == nil {
			log.Println("client:"+r.RemoteAddr, " req:", string(data))
		} else {
			log.Println("DumpRequest:", err.Error())
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func ReadBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get body
		if r.Method == "POST" || r.Method == "PUT" {
			v, ok := api_struct_map[KeyFromRequestURI(r.URL.Path)]

			if ok && v != nil {
				b, err := ioutil.ReadAll(r.Body)
				r.Body.Close()
				if err != nil {
					log.Println(err)
					funcErrResponse(w, err)
					return
				}

				if reflect.Map == reflect.TypeOf(v).Kind() {
					m := map[string]interface{}{}
					err = json.Unmarshal(b, &m)
					if err != nil {
						log.Println(err)
						funcErrResponse(w, err)
						return
					}
					r = r.WithContext(context.WithValue(r.Context(), "param", m))
				} else {
					m_ptr := reflect.New(reflect.TypeOf(v)).Interface()

					initIntField(m_ptr)

					// m := reflect.Zero( reflect.TypeOf( v ) ).Interface()
					err = json.Unmarshal(b, m_ptr)
					if err != nil {
						log.Println(err)
						funcErrResponse(w, err)
						return
					}

					err = verifyParams(m_ptr)
					if err != nil {
						funcErrResponse(w, err)
						return
					}

					r = r.WithContext(context.WithValue(r.Context(), "param", m_ptr))
				}

			}
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func GetAPIDoc() string {
	return api_doc
}

func RegisterAPI(r *mux.Router, apis []API) {

	for _, api := range apis {
		mappingApiStruct(api.Path, api.Params, api.Method)
		r.HandleFunc(api.Path, api.Func).Methods(api.Method)

		api_doc += fmt.Sprintf("    %s\n", KeyFromRequestURI(api.Path))
		api_doc += fmt.Sprintf("\tpath: '%s'\n", api.Path)
		api_doc += fmt.Sprintf("\tmethod:%s\n", api.Method)
		api_doc += fmt.Sprintf("\tbody:%+v", api.Params)
		if api.Params != nil {
			api_doc += fmt.Sprintf(", (%s)\n", "please lower-case json key")
		} else {
			api_doc += fmt.Sprintf("\n")
		}
		api_doc += fmt.Sprintf("\texample: %s\n", api.Example)
		api_doc += fmt.Sprintf("\tdesc: %s\n", api.Desc)
	}
	log.Println(strings.Replace(api_doc, "{DOMAIN}", domain, -1))
}
