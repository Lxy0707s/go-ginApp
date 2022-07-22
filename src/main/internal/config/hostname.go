package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	hostName string
)

func InitHostname() {
	if hostName == "" {
		hostName, _ = GetCurrentComputerHostName("/allconf/hostname.conf")
	}
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// GetCurrentComputerHostName 获取当前机器的HostName
func GetCurrentComputerHostName(path string) (string, error) {
	filePath := "/allconf/hostname.conf"
	if notExist := CheckNotExist(filePath); notExist {
		name, _ := os.Hostname()
		return name, nil
	}
	if path != "" {
		filePath = path
	}
	hostname, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("zeus program close read hostname error:", err)
		hostName, err := os.Hostname()
		return hostName, err
	}
	// 去除首尾空格
	hostName := strings.Trim(string(hostname[:]), " ")
	hostName = strings.Replace(hostName, "\n", "", -1)
	hostName = strings.Replace(hostName, "hostname=", "", -1)
	return hostName, nil
}
