package httpClient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"go-ginApp/src/main/pkg/common/httpserver/httpServer"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	// ClientOption is the Option of HTTP Client
	ClientOption struct {
		Name        string            `json:"name"`
		URL         string            `json:"url"`
		Retry       int               `json:"retry"`
		Timeout     int               `json:"timeout"`
		MaxConn     int               `json:"max_conn"`
		User        string            `json:"-"`
		Password    string            `json:"-"`
		ExtraHeader map[string]string `json:"header"`
		HashKey     string            `json:"hash_key"`
	}

	// Client HTTP Client
	Client struct {
		name string
		opt  *ClientOption
		tr   *http.Transport
	}
)

func fillOption(opt *ClientOption) *ClientOption {
	if opt.Name == "" {
		opt.Name = "client"
	}
	if opt.Retry < 1 {
		opt.Retry = 3
	}
	if opt.Timeout < 1 {
		opt.Timeout = 10
	}
	if opt.MaxConn < 1 {
		opt.MaxConn = 100
	}
	return opt
}

// NewClient init a new Client
func NewClient(opt *ClientOption) *Client {
	if opt.URL == "" {
		return nil
	}
	opt = fillOption(opt)
	return &Client{
		name: opt.Name,
		opt:  opt,
		tr: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          opt.MaxConn,
			MaxIdleConnsPerHost:   opt.MaxConn,
			IdleConnTimeout:       60 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// ChangeURL change the url of Client
func (c *Client) ChangeURL(url string) {
	c.opt.URL = url
}

// SendGET send get request
func (c *Client) SendGET(paramsMap map[string]string) (respondBody []byte, statusCode int, err error) {
	var params string
	if len(paramsMap) > 0 {
		values := url.Values{}
		for key, value := range paramsMap {
			values.Set(key, value)
		}
		if strings.Contains(c.opt.URL, "?") {
			params = "&" + values.Encode()
		} else {
			params = "?" + values.Encode()
		}
	}
	return c.send("GET", c.opt.URL+params, bytes.NewReader([]byte("")))
}

// SendPOST send body with client
// it retries opt.Retry times
// if all times are failed, return error
func (c *Client) SendPOST(body []byte) error {
	var err error
	for i := 0; i <= c.opt.Retry; i++ {
		bodyReader := bytes.NewReader(body)
		if _, _, err = c.sendOnce("POST", c.opt.URL, bodyReader); err == nil {
			break
		}
	}
	return err
}

func (c *Client) send(sendType string, url string, body io.Reader) (respondBody []byte, statusCode int, err error) {
	for i := 0; i <= c.opt.Retry; i++ {
		if respondBody, statusCode, err = c.sendOnce(sendType, url, body); err == nil {
			break
		}
	}
	return
}

func (c *Client) sendOnce(sendType string, url string, body io.Reader) (respondBody []byte, statusCode int, err error) {
	req, err := http.NewRequest(sendType, url, body)
	if err != nil {
		return nil, 0, err
	}
	if c.opt.User != "" {
		req.SetBasicAuth(c.opt.User, c.opt.Password)
	}
	if len(c.opt.ExtraHeader) > 0 {
		for k, v := range c.opt.ExtraHeader {
			req.Header.Set(k, v)
		}
	}
	if c.opt.HashKey != "" {
		for k, v := range httpServer.BuildHashHeader(c.opt.HashKey) {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{
		Timeout:   time.Second * time.Duration(c.opt.Timeout),
		Transport: c.tr,
	}
	resp, err := client.Do(req)

	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("bad status :%d", resp.StatusCode)
	}
	statusCode = resp.StatusCode
	respondBody, err = ioutil.ReadAll(resp.Body)
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return
}
