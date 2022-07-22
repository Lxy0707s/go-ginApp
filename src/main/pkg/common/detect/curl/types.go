package main

import "net/http"

type CurlParams struct {
	Http      *http.Server
	Url       string          //目标url
	HttpType  string          //请求方式，默认get
	Timeout   int             // 超时时间 秒
	Transport *http.Transport // 请求配置
	Headers   *http.Header    // 头部信息
}
