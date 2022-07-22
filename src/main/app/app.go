package app

import (
	"go-ginApp/src/main/internal/config"
	"go-ginApp/src/main/internal/middleware"
	"go-ginApp/src/main/internal/servers"
	"go-ginApp/src/main/pkg/utils/dbtool"
	"log"
)

var ins *servers.HTTPServer

func InitApp(version, customConfigName string) {
	in := middleware.NewInstance()
	in.Demo(nil)
	runModel, configName := FlagInit(version)
	if customConfigName != "" {
		configName = customConfigName
	}
	// 读取配置文件
	cfg := config.InitConfig(runModel, configName)
	config.SetConfig(cfg)
	// 加载数据库
	dbtool.Setup(config.AppConfig.Database, true)
	// 开启http服务
	ins = servers.HttpInstance()
	runServer()
}

func runServer() {
	go func() {
		err := ins.Start()
		if err != nil {
			log.Println(err)
		}
	}()
}
