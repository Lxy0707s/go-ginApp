package config

import (
	"github.com/spf13/viper"
	"go-ginApp/src/main/pkg/utils/base_struct"
	"go-ginApp/src/main/pkg/utils/dbtool"
	"go-ginApp/src/main/pkg/utils/logtool"
	"sync"
)

var (
	CfgInfo Config
	lock    sync.RWMutex

	Viper     *viper.Viper
	AppLog    logtool.Logger
	AppConfig = &Config{}
)

type Config struct {
	AppName   string
	Version   string
	Author    string
	Database  dbtool.Option
	System    base_struct.SystemConfig
	Server    base_struct.ServerConfig
	ApiTokens map[string]base_struct.ApiTokenMap
	JwtConfig base_struct.JwtConfig
	//...
}
