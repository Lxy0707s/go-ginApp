package cmd

import (
	"fmt"
	"go-ginApp/src/main/pkg/utils/logtool"
	"os/exec"
	"strings"
	"sync"
)

var (
	cmdInstance *CoreCmd
	once        sync.Once
)

// RunCmd 定义接口规范和核心服务
type (
	CommandSplit func(cmdStr string) string
	OutSplit     func(out string) interface{}
	CommandList  interface {
		GenerateCmd(cmdStr string, t string, c ...CommandSplit) string
		RunCmd(cmdStr string, os ...OutSplit) interface{}
	}
	CoreCmd struct {
		log  logtool.Logger
		lock sync.RWMutex
	}
)

// NewCmdInstance 初始化命令执行服务
func NewCmdInstance(server string) *CoreCmd {
	if cmdInstance == nil {
		once.Do(func() {
			cmdInstance = &CoreCmd{
				log: logtool.NewSugar("cmd-server", false),
			}
		})
	}
	return cmdInstance
}

// GenerateCmd 命令初始化
func (c *CoreCmd) GenerateCmd(cmdString string, _type string, cmdSplit ...CommandSplit) string {
	var command = ""
	switch _type {
	case SYSTEMCTL: // 服务管理 例子
		if cmdString != "" {
			command = cmdString
			for _, splitFunc := range cmdSplit {
				command = splitFunc(cmdString)
			}
		}
		return command
	case SUPERVISOR: // 服务管理
	case CURL: // 文件下载
	case PING: // 网络连通性测试
	case UDP: // 端口扫描
	case TCP: // 端口扫描，tcp连接检测
	case MTR: // mtr检测
	case TRACEROUTE: // traceroute跳检测
	case DIG: // 域名解析
	case NSLOOKUP: // 域名解析

	}
	return command
}

// RunCmd 命令执行，返回想要的结果
func (c *CoreCmd) RunCmd(cmdStr string, outSplit ...OutSplit) (string, error) {
	process := exec.Command("/bin/sh", "-c", cmdStr)
	out, err := process.Output()
	if err != nil {
		return "", err
	}
	data := strings.Trim(string(out), "\n")
	var cmdResult = ""
	fmt.Println("执行指令", "command", cmdStr)
	if outSplit != nil {
		for _, cmdOut := range outSplit {
			cmdResult = cmdOut(data).(string)
		}
	}
	return cmdResult, nil
}

type Demo struct {
}

func (d *Demo) Test() {
	// runCmd()
	cmdStr := "systemctl status docker | grep Active"
	cmd := NewCmdInstance("test")
	cmdStr = d.cmdSplit(cmdStr) //自定义命令拼接
	cmdStr = cmd.GenerateCmd(cmdStr, "systemctl")
	fmt.Println(cmd.RunCmd(cmdStr))
}

func (d *Demo) cmdSplit(cmdString string) (cmd string) {
	cmd = cmdString

	return cmd
}
