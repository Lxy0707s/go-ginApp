package detect

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func GetOutBoundIP() (ip string, port string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.String(), ":")[0]
	port = strings.Split(localAddr.String(), ":")[1]
	return
}

func GetIPV4() string {
	resp, err := http.Get("https://ipv4.netarm.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}
