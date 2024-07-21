package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wenxuan7/solution/external"
)

func setup() {
	// 连接mysql数据库
	external.Mysql()
	// 连接redis
	external.Redis()
}

func main() {
	setup()
	// gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	omsGroup := router.Group("/oms", func(c *gin.Context) {

	})
	// 配置
	settingsController(omsGroup)
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
