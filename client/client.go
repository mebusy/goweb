package client

import (
	"log"
	// "crypto/sha256"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
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

	nConcurrency := 128
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		MaxConnsPerHost: nConcurrency,
		//MaxIdleConns:        nConcurrency,
		//MaxIdleConnsPerHost: nConcurrency,
	}
	if proxyURL != nil {
		tr.Proxy = http.ProxyURL(proxyURL)
	}
	http_client = &http.Client{Timeout: 15 * time.Second, Transport: tr}
}

func Client() *http.Client {
	return http_client
}

func readResponse(res *http.Response, err error) (string, error) {

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
		return "", errors.New(str_body)
	}

	return str_body, nil

}

func Get(url string) (string, error) {
	res, err := http_client.Get(url)
	return readResponse(res, err)
}

func Post(url string, contentType string, data io.Reader) (string, error) {
	res, err := http_client.Post(url, contentType, data)
	return readResponse(res, err)
}
func PostWithCookie(url string, contentType string, data io.Reader, cookie string ) (string, error) {
    req, err := http.NewRequest("POST", url , data )
    if err != nil {
        return readResponse( nil, err)
    }
    if contentType != "" {
        req.Header.Set( "Content-Type" , contentType )
    }
    if cookie != "" {
        req.Header.Set( "Cookie", cookie )
    }
    res, err := http_client.Do(req)
	return readResponse(res, err)
}


