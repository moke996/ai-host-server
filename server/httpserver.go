package server

import (
	"ai-host/controller"
	"ai-host/global"
	"ai-host/middleware"
	"github.com/gin-gonic/gin"
)

func NewHttp() {
	r := gin.Default()
	r.Use(middleware.CorsMiddleware())
	r.GET("/api/v1/host/versions", controller.GetVersion)
	r.POST("/api/v1/host/start", controller.Start)
	r.POST("/api/v1/host/run", controller.Run)
	r.POST("/api/v1/host/save", controller.Save)
	r.GET("/api/v1/host/histories", controller.GetHistoryList)
	r.GET("/api/v1/host/history/:id", controller.GetHistoryById)
	r.POST("/api/v1/host/history/:id", controller.DeleteHistoryById)
	_ = r.Run(global.Config.Common.HttpAddress)
}
