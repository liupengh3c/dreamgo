package main

import (
	"female/lib/tools"
	"female/router"

	"github.com/BurntSushi/toml"

	"github.com/gin-gonic/gin"
)

type App struct {
	AppName    string `toml:"app_name"`
	RunMode    string `toml:"run_mode"`
	HTTPListen string `toml:"http_listen"`
}

//入口函数，所有http 请求全部请求到这里，之后根据路由进行分发
func main() {
	// engine := gin.Default()
	config := new(App)
	toml.DecodeFile(tools.GetCurrentDirectory()+"/conf/app.toml", config)
	engine := gin.New()
	engine.Use(gin.Recovery())
	router.Register(engine)
	engine.Run(config.HTTPListen)
}
