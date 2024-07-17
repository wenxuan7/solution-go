package main

import "github.com/wenxuan7/solution/link"

func main() {
	// 连接mysql数据库
	link.Mysql()
	// 连接redis
	link.Redis()
	// 本地缓存 bigCache
}
