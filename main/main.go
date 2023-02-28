package main

import (
	"go.uber.org/zap"
	"solution-go/cache"
	"solution-go/db"
)

func main() {
	var (
		err    error
		logger *zap.Logger
	)

	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	db.ConnectMysql()
	cache.ConnectRedis()

	err = cache.Rdb.Set(cache.Ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := cache.Rdb.Get(cache.Ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	sugar.Infof("key -> %s", val)
}
