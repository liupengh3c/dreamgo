package router

import (
	"female/controllers"
	"female/middlewares"

	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {
	engine.Use(middlewares.Print)
	engine.POST("female/upload", controllers.Upload)
	engine.POST("female/sample", controllers.Sample)
	engine.GET("female/test", controllers.Test)
}
