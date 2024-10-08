package main

import (
	"ai-dating/controller"
	"ai-dating/global"
	"ai-dating/server"
	"time"
)

const GracefulStopInterval = 15 * time.Second

func main() {
	// 加载配置
	global.LoadConfig()
	// 初始化依赖
	global.Init()
	// 加载数据
	controller.InitData()
	// 启动http服务
	server.NewHttp()
}
