package config

import (
	"github.com/spf13/viper"
	"go-ginApp/src/main/pkg/utils/logtool"
	"log"
	"time"
)

func InitConfig(runModel, configName string) *Config {
	AppLog = logtool.NewSugar("demo", true)
	log.Println("Now Run Model: ", runModel)
	// 初始化加载目标配置文件
	CfgInfo = readFileConfig(configName)
	go Loop(configName)
	return &CfgInfo
}

func SetConfig(config *Config) {
	lock.Lock()
	defer lock.Unlock()
	AppConfig = config
}

func Loop(configName string) {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			log.Println("Get Config Information .......")
			lock.Lock()
			CfgInfo = readFileConfig(configName)
			SetConfig(&CfgInfo)
			lock.Unlock()
		}
	}()
	select {}
}

func readFileConfig(configName string) Config {
	log.Println("read config file starting...")
	var configuration Config
	viper.SetConfigFile(configName)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	unmarshalFileConfig(&configuration)
	return configuration
}

func unmarshalFileConfig(configuration *Config) {
	err := viper.Unmarshal(configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
