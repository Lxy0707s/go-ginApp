package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-ginApp/src/main/pkg/utils/logtool"
	"go-ginApp/src/main/pkg/utils/myfile"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

type job struct {
	log   logtool.Logger
	param *CurlParams
}

var instance *job

func NewInstance() *job {
	if instance == nil {
		instance = &job{
			log: logtool.NewSugar("hcurl-job", false),
		}
	}
	return instance
}

func (j *job) GenerateExec() {
	// 1 -x xxxx:port  代理信息
	u, _ := url.Parse("127.0.0.1:80")
	// 2 --resolve 强制解析指定
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	originAddr := "www.baidu.com:443"
	resolvesAddr := "163.177.151.110:443"
	flag := false
	// or create your own transport, there's an example on godoc.
	tr := &http.Transport{
		MaxIdleConns:    10,
		MaxConnsPerHost: 10,
		IdleConnTimeout: time.Duration(10) * time.Second,
		Proxy:           http.ProxyURL(u),
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == originAddr {
				addr = resolvesAddr
			}
			return dialer.DialContext(ctx, network, addr)
		},
		// -k 是否跳过ssl检测
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: flag,
		},
	}
	// 3 头部信息
	h := &http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("User-Agent", "curl/7.29.0")
	// 4 -L 重定向
	redirectUrl := ""
	refererUrl := ""
	if true && refererUrl != "" {
		h.Set("Location", redirectUrl)
	}
	// 5 --referer 跳转查询，表示你是从哪里跳转过来的
	if true && refererUrl != "" {
		h.Set("Referer", refererUrl)
	}
	// 6 --cookie设置
	cookie := ""
	if true {
		h.Set("Cookie", cookie)
	}
	j.param = &CurlParams{
		Url:       "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png", // 目标url
		HttpType:  "get",                                                                 // 请求方式
		Timeout:   10,                                                                    // 超时
		Transport: tr,
		Headers:   h,
	}
}

func (j *job) exec() {
	durl := j.param.Url
	///// 指定ip【代理】的情况 -x http
	t := j.param.Transport
	client := &http.Client{
		Transport: t,
		Timeout:   time.Duration(j.param.Timeout) * time.Second,
	}
	var start, connect, dns, tlsHandshake time.Time
	var dNSDoneTime, tlsHandshakeTime, connectTime, firstRespByteTime int64

	req, _ := http.NewRequest("GET", durl, nil)
	// 创建客户端请求跟踪
	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) {
			dns = time.Now()
		},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dNSDoneTime = time.Since(dns).Milliseconds()
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tlsHandshakeTime = time.Since(tlsHandshake).Milliseconds()
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			connectTime = time.Since(connect).Milliseconds()
		},

		GotFirstResponseByte: func() {
			firstRespByteTime = time.Since(start).Milliseconds()
		},
	}
	// 请求追踪
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	// 新增头部
	req.Header.Add("Authorization", "token")
	req.Header.Add("Content-Type", "application/json")
	//resp, err := http.DefaultTransport.RoundTrip(req)
	start = time.Now()
	resp, err := client.Do(req)
	if err != nil {
		j.log.Info("error", ":", err)
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	myfile.WriteJSON(string(data), "data.png")
	totalTime := time.Since(start) // 总耗时
	fileSize := resp.ContentLength // 获取响应内容大小
	statusCode := resp.StatusCode  // 响应状态码
	defer resp.Body.Close()
	tt := totalTime.Seconds()
	if tt == 0 {
		tt = 1
	}
	speed := float64(fileSize) / tt
	downloadSpeed := bytesToSize(speed, 1.0) //下载速度

	fmt.Println("------------start----------")
	fmt.Printf("Http status: %v\n", statusCode)                         // statusCode
	fmt.Printf("DNS Done(ms): %v\n", dNSDoneTime)                       // DNS解析耗时
	fmt.Printf("Connect time(ms): %v\n", connectTime)                   // 建联时间
	fmt.Printf("TLS Handshake(ms): %v\n", tlsHandshakeTime)             // 握手耗时
	fmt.Printf("From start to first byte(ms): %v\n", firstRespByteTime) // 首包
	fmt.Printf("Totle time(ms): %v\n", totalTime.Microseconds())        // totalTime
	fmt.Printf("Totle time(s): %f\n", totalTime.Seconds())              // totalTime
	fmt.Printf("File size(Bytes): %v\n", fileSize)                      // 文件大小
	fmt.Printf("Download speed(KB/s): %f\n", downloadSpeed)             // downloadSpeed
	fmt.Println("-----------end-----------")
}

func bytesToSize(length, size float64) float64 {
	var k = 1024.0 // or 1024
	if length == 0 {
		return 0
	}
	k = math.Pow(k, size)
	//var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	// i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	//r := float64(length) / math.Pow(float64(k), i)
	//strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
	r := length / k
	return r
}
