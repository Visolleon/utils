package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//FROM: https://github.com/gopherchina/website/blob/master/models/http.go
var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.120 Safari/537.36"
var (
	dialTimeout    = 5 * time.Second   //flag.Duration("dial_timeout", 10*time.Second, "Timeout for dialing an HTTP connection.")
	requestTimeout = 120 * time.Second //flag.Duration("request_timeout", 20*time.Second, "Time out for roundtripping an HTTP request.")
)

func timeoutDial(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, dialTimeout)
}

type transport struct {
	t http.Transport
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	timer := time.AfterFunc(requestTimeout, func() {
		t.t.CancelRequest(req)
		//beego.Warn("Canceled request for %s", req.URL)
	})
	defer timer.Stop()
	resp, err := t.t.RoundTrip(req)
	return resp, err
}

var (
	httpTransport = &transport{t: http.Transport{Dial: timeoutDial, ResponseHeaderTimeout: requestTimeout / 2}}
	httpClient    = &http.Client{Transport: httpTransport}
)

// GetJSON GET Get JSON
func GetJSON(url string, v interface{}) error {
	return doJSONReq("GET", url, nil, v)
}

// PostJSON POST Get JSON
func PostJSON(url string, data url.Values, v interface{}) error {
	return doJSONReq("POST", url, data, v)
}

// Request Get JSON FROM HTTP URL
func Request(method string, url string, data url.Values) (int, []byte, error) {
	var body io.Reader
	isPost := false
	if method == "POST" {
		isPost = true
		if data != nil && len(data) > 0 {
			body = ioutil.NopCloser(strings.NewReader(data.Encode()))
		}
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	if isPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	bodyByte, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Response Body:", string(bodyByte))
	return resp.StatusCode, bodyByte, err
}

func doJSONReq(method string, url string, data url.Values, v interface{}) error {
	statusCode, body, err := Request(method, url, data)
	if statusCode == 200 && err == nil {
		err = json.Unmarshal(body, v)
		if _, ok := err.(*json.SyntaxError); ok {
			return errors.New("JSON syntax error at " + url)
		} else if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Status: %d", statusCode)
}
