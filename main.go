package main

import "go-ginApp/src/main/app"

const (
	AppName    = "go-ginApp"
	AppVersion = "1.5.2"
	Author     = "xuanyu.li"
	BuildTime  = "2022-01-25 11:00:00"
)

func main() {
	app.InitApp(AppVersion, "")
}
