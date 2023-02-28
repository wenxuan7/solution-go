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
		val    string
	)

	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	db.ConnectMysql()
	cache.ConnectRedis()

	if err = cache.Rdb.Set(cache.Ctx, "key", "value", 0).Err(); err != nil {
		panic(err)
	}

	if val, err = cache.Rdb.Get(cache.Ctx, "key").Result(); err != nil {
		panic(err)
	}

	sugar.Infof("key -> %s", val)
}
