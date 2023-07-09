package main

import (
	"MYBLOGWEB/common"
	"MYBLOGWEB/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	//初始化数据库
	db := common.InitDB()
	//延迟关闭数据库
	defer db.Close()
	router := gin.Default()
	// 使用 CORS 中间件
	router.Use(cors.Default())
	//打开路由
	router = routers.CollectRouter(router)
	//设置监听端口
	router.Run(":8080")
}
