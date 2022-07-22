package httpServer

import (
	"crypto/tls"
	"errors"
	"github.com/tidwall/gjson"
	"go-ginApp/src/main/pkg/utils/prof"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

/*
requestServer 提供请求数据的方法，
同时统计请求的各个接口的请求次数，QPS，请求错误次数，请求处理最大耗时，平均处理耗时。
*/

type (
	API struct {
		API      string            `json:"api,omitempty"`     // 地址
		Method   string            `json:"method,omitempty"`  // POST GET
		Token    string            `json:"token,omitempty"`   // 凭证
		Params   string            `json:"params,omitempty"`  // 参数
		Timeout  int64             `json:"timeout,omitempty"` // 超时时间
		Retry    int64             `json:"retry,omitempty"`   // 重试次数
		Header   map[string]string `json:"header,omitempty"`  // 头
		Interval int64             `json:"interval,omitempty"`
	}
	requestServer struct {
		lock    sync.RWMutex
		reqCnt  map[string]*prof.CountQPS
		errCnt  map[string]*prof.CountBase
		reqTime map[string]*prof.AverageTimer
	}
	RequestHandler interface {
		Handler(api API) ([]byte, error)
	}
)

var (
	once    sync.Once
	request *requestServer
)

func GetInstance() *requestServer {
	if request == nil {
		once.Do(func() {
			request = &requestServer{
				reqCnt:  make(map[string]*prof.CountQPS),
				errCnt:  make(map[string]*prof.CountBase),
				reqTime: make(map[string]*prof.AverageTimer),
			}
		})
	}
	return request
}
func (r *requestServer) ProfMap() map[string][]interface{} {
	profs := make(map[string][]interface{})
	for name, p := range r.reqCnt {
		profs[name] = append(profs[name], p)
	}
	for name, p := range r.errCnt {
		profs[name] = append(profs[name], p)
	}
	for name, p := range r.reqTime {
		profs[name] = append(profs[name], p)
	}
	return profs
}
func (r *requestServer) Handler(api API) ([]byte, error) {
	// 尝试创建计算器
	r.lock.Lock()
	if _, ok := r.reqCnt[api.API]; !ok {
		r.reqCnt[api.API] = prof.NewCountQPS("req")
		r.reqTime[api.API] = prof.NewAverageTimer("time", 0)
		r.errCnt[api.API] = prof.NewCountBase("error")
	}
	r.lock.Unlock()
	// 统计处理耗时
	start := time.Now()
	defer func() {
		r.reqTime[api.API].Set(time.Since(start).Milliseconds())
	}()

	// 统计请求次数
	r.reqCnt[api.API].Incr()

	// 构建请求客户端 默认跳过tls验证
	client := &http.Client{
		Timeout: time.Second * time.Duration(api.Timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest(api.Method, api.API, strings.NewReader(api.Params))
	if err != nil {
		r.errCnt[api.API].Incr()
		return nil, err
	}
	// 尝试加入token
	if api.Token != "" {
		req.Header.Add("Authorization", api.Token)
	}
	// 默认加入Content-Type
	req.Header.Add("Content-Type", "application/json")
	// 尝试加入header
	if api.Header != nil {
		for key := range api.Header {
			req.Header.Add(key, api.Header[key])
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		r.errCnt[api.API].Incr()
		return nil, err
	}
	defer resp.Body.Close()

	// 读出内容 默认对内容进行校验 code 不为0 或者 200；data 为空 默认请求失败，
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.errCnt[api.API].Incr()
		return nil, err
	}

	jsonResult := gjson.ParseBytes(body)
	dataStr := jsonResult.Get("data").String()
	code := jsonResult.Get("code").String()
	if code != "0" && code != "200" {
		r.errCnt[api.API].Incr()
		return body, errors.New("request fail code " + code)
	}
	if dataStr == "{}" {
		r.errCnt[api.API].Incr()
		return body, errors.New("request fail data is nil ")
	}
	return body, nil
}
func Handler(api API) ([]byte, error) {
	if request == nil {
		GetInstance()
	}
	return request.Handler(api)
}
