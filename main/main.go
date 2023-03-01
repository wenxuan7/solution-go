package main

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"solution-go/cache"
	"solution-go/db"
)

func main() {
	var (
		err       error
		logger    *zap.Logger
		val       string
		statusCmd *redis.StatusCmd
	)

	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	db.ConnectMysql()
	cache.ConnectRedis()

	statusCmd = cache.Rdb.Set(cache.Ctx, "key", "value", 0)
	if err = statusCmd.Err(); err != nil {
		panic(err)
	}

	if val, err = cache.Rdb.Get(cache.Ctx, "key").Result(); err != nil {
		panic(err)
	}

	sugar.Infof("key -> %s", val)
}
