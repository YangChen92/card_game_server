package main

import (
	"game_server/global"
	"game_server/handler"
	"game_server/initialize"

	"github.com/aceld/zinx/znet"
)

func main() {
	// 1. 初始化配置
	initialize.InitConfig()

	// 2. 初始化数据库
	initialize.InitMySQL()

	// 3. 初始化Redis
	initialize.InitRedis()

	// 4. 创建服务器
	s := znet.NewServer()

	// // 5. 添加全局中间件
	// var tokenAuthMiddleware = &middleware.TokenAuthMiddleware{}
	// s.Use(tokenAuthMiddleware.RouterHandler)

	// 6. 注册路由
	s.AddRouter(global.MSG_REGISTER, &handler.UserHandler{})
	s.AddRouter(global.MSG_LOGIN, &handler.UserHandler{})

	// 7. 启动服务
	s.Serve()
}
