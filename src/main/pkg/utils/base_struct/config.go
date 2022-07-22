package base_struct

import "time"

type (
	SystemConfig struct {
		RuntimeRootPath string
		Debug           bool
		HostNamePath    string
		Machine         MachineList
	}
	ServerConfig struct {
		RunMod            string
		Addr              string
		GRPCAddr          string
		WriteTimeout      time.Duration
		ReadTimeout       time.Duration
		CertFile          string
		KeyFile           string
		APIFrequencyLimit int
	}
	ApiTokenMap struct {
		MaxFileSize     int64
		AppID           string
		AppToken        string
		SourceApp       string
		AimApp          string
		QueryTimeLimit  int
		QueryCountLimit int
	}
	JwtConfig struct {
		Switch bool
		//刷新用户相关请求参数（from sso)
		AppId            string
		AppSecret        string
		AppSalt          string
		JwtUrl           string
		TokenParseSwitch bool
	}
)

type MachineList struct {
	Master []string //主机列表
	Backup []string //备机列表
}
