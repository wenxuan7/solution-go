package main

import (
	"github.com/solution-go/cache"
	"github.com/solution-go/db"
	"github.com/solution-go/log"
	"github.com/solution-go/web"
)

func main() {
	// 初始化log
	log.InitSugaredLogger()
	// 连接mysql
	db.ConnectMysql()
	// 连接redis
	cache.ConnectRedis()
	// 启动gin
	web.StartGin()
}
