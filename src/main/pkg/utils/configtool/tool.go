package configtool

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"unicode/utf8"
)

// ErrorFileCheck 进程文件无法找到的错误
var ErrorFileCheck error = errors.New("ExecutorSockPath file is not exist")

/**
* @Author: Shaun Chen
* @Date: 2019-06-24 15:56
 */

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readPortFile(filePath string) (int64, error) {
	var portInt64 int64
	fileExist, _ := pathExists(filePath)
	if fileExist {
		portBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return 0, err
		}
		portStr := string(portBytes)
		if portStr != "" {
			portInt64, _ = strconv.ParseInt(portStr, 10, 64)
		}
	}
	return portInt64, nil
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
