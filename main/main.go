package main

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/solution-go/cache"
	"github.com/solution-go/data"
	"github.com/solution-go/db"
	"go.uber.org/zap"
)

func main() {
	var (
		err       error
		logger    *zap.Logger
		val       string
		statusCmd *redis.StatusCmd
		tm        *data.TradeMain
		bs        []byte
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

	tm = &data.TradeMain{}
	db.MysqlDB.First(tm)
	if bs, err = jsoniter.Marshal(tm); err != nil {
		sugar.Fatalln(err)
	}

	sugar.Info(string(bs))
}
