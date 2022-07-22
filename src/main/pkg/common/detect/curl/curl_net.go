package main

import (
	"fmt"
	"strings"
)

var durl = "https://www.cnipa.gov.cn/picture/0/2006101108210115064.png"
var CurlExecParams = []string{
	HParams,
	ResolveParams,
	RefererParams,
	CookieParams,
	LParams,
	HttpParams,
}

const (
	HParams       = "-H"
	HttpParams    = "-X"
	ResolveParams = "--resolve"
	RefererParams = "--referer"
	CookieParams  = "--cookie"
	LParams       = "-L"
)

type dJob struct {
	DUrl       *string
	resolveUrl *string
	refererStr *string
}

func main() {
	/*	job := NewInstance()
		job.GenerateExec()
		job.exec()
	*/
	c := "256K.dat -H 'Host:jiankongtest.baishancdnx.cn' -X 'POST' --resolve '123'"
	//c := "https://123.com -H 'Host:127.0.0.1' --referer 'Referer:xxx' --resolve 'www.baidu.com:443:127.0.0.1:443' -L 'ccc'"
	var r = make(map[string][]string)
	a := getExecElement(c, r)
	if a == "" {
		// 兼容目标url在字符串首部的情况 "https://123.com -H 'Host:127.0.0.1' --referer 'Referer:xxx' --resolve 'www.baidu.com:443:127.0.0.1:443' -L 'ccc'"
		a = strings.Split(c, " ")[0]
	}
	if _, ok := r[HttpParams]; ok && len(r[HttpParams]) == 1 {
		fmt.Println("9999999999", r[HttpParams][0])
	}
	fmt.Println(r)
	fmt.Println("---------")
	fmt.Println(a)
}

// 处理带参数的URL :256K.dat -H 'Host:jiankongtest.baishancdnx.cn'
func getExecElement(s string, res map[string][]string) string {
	var result = ""
	for _, p := range CurlExecParams {
		result = getElement(s, p, p, res)
		/*switch p {
		case ResolveParams: //只会有一个
			result = getElement(s, ResolveParams, ResolveParams, res)
		case HParams:
			result = getElement(s, HParams, HParams, res)
		case LParams:
			result = getElement(s, LParams, LParams, res)
		case RefererParams:
			result = getElement(s, RefererParams, RefererParams, res)
		case CookieParams:
			result = getElement(s, CookieParams, CookieParams, res)
		case HttpParams:
			result = getElement(s, HttpParams, HttpParams, res)
		}*/
	}
	return result
}

// 迭代获取同类参数
func getElement(str, symbol, name string, res map[string][]string) string {
	index := strings.Index(str, symbol)
	var result = ""
	if index != -1 {
		newStr := str[index:]
		newStrTrim := strings.Trim(strings.Split(newStr, symbol)[1], " ")
		newStrs := strings.Split(newStrTrim, " ")
		if newStrs[0] == "" {
			return newStrTrim
		}
		if symbol == HttpParams {
			fmt.Println(str, newStr, newStrs)
		}
		res[name] = append(res[name], newStrs[0])
		if len(newStrs) > 1 {
			getExecElement(newStrs[1], res)
			result = newStrs[1]
		}
	}
	return result
}
