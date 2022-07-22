package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
)

// Flags is a simple flag interrupter to print value and load correct config file
func Flags(version string, buildTime string, cfg interface{}) {
	c := flag.Bool("c", false, "show default config")
	v := flag.Bool("v", false, "show version")
	vt := flag.Bool("vt", false, "show version and built time")
	flag.Parse()
	if *c {
		b, _ := json.MarshalIndent(cfg, "", "\t")
		fmt.Println(string(b))
		os.Exit(0)
	}
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}
	if *vt {
		fmt.Println("version : " + version)
		fmt.Println("build : " + buildTime)
		fmt.Println("go : " + runtime.Version())
		fmt.Println("os : " + runtime.GOOS + "/" + runtime.GOARCH)
		os.Exit(0)
	}
}

func FlagInit(version string) (string, string) {
	v := flag.Bool("v", false, "show version")
	var runModel, configName string
	flag.StringVar(&runModel, "runModel", "local", "set run model")
	flag.StringVar(&configName, "configName", "config.yml", "Sets the configuration path but does not require the final format")
	vt := flag.Bool("vt", false, "show version and built time")
	flag.Parse()
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}
	if *vt {
		fmt.Println("version : " + version)
		fmt.Println("run model : " + runModel)
		fmt.Println("config name : " + configName)
		fmt.Println("go : " + runtime.Version())
		fmt.Println("os : " + runtime.GOOS + "/" + runtime.GOARCH)
		os.Exit(0)
	}
	//listFlag, startFlag, stopFlag, lsFlag := RegisterFlag()
	//FlagHandle(listFlag, startFlag, stopFlag, lsFlag)
	return runModel, configName
}

// RegisterFlag 添加默认命令行参数
func RegisterFlag() (listFlag *bool, startFlag *string, stopFlag *string, lsFlag *bool) {
	listFlag = flag.Bool("list", false, "list services")
	lsFlag = flag.Bool("ls", false, "list services")
	startFlag = flag.String("start", "", "start services")
	stopFlag = flag.String("stop", "", "stop services")
	return
}
