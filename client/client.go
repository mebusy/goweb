package client

import (
    "log"
	// "crypto/sha256"
	"crypto/tls"
	"net/http"
    "net/url"
	"os"
    "time"
    "io/ioutil"
    "errors"
    "io"
)

var http_client *http.Client

func init() {
    var proxyURL *url.URL = nil
	http_proxy := os.Getenv("HTTP_PROXY")
	if http_proxy != "" {
		p, err := url.Parse(http_proxy)
		if err != nil {
			log.Println(err)
		} else {
			proxyURL = p
			log.Println("proxyURL ", proxyURL)
		}
	}
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    if proxyURL != nil {
        tr.Proxy = http.ProxyURL(proxyURL)
    }
    http_client = &http.Client{Timeout: 15 * time.Second, Transport: tr}
}


func Client() *http.Client {
    return http_client 
}


func readResponse( res *http.Response , err error  )  ( string, error  ) {

    if err != nil {
        log.Println(err)
        return "", err
    }

    b, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        log.Println(err)
        return "", err
    }

    str_body := string(b)

    if res.StatusCode >= 200 && res.StatusCode <= 299 {
    } else {
        return ""  ,  errors.New( str_body )  
    }
    
    return str_body , nil 

}

func Get( url string ) ( string , error  ) {
    res, err := http_client.Get(url)
    return readResponse( res, err )
}



func Post( url string , contentType string  ,data io.Reader ) ( string , error  ) {
    res, err := http_client.Post( url, contentType  , data  )
    return readResponse( res, err )
}

