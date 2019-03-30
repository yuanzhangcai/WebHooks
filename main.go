package main

import (
	_ "WebHooks/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	logSetting := beego.AppConfig.String("logSetting")
	beego.SetLogger("file", logSetting)

	beego.Run()
}
