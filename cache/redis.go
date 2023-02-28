package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"os"
)

var (
	Ctx = context.Background()
	Rdb *redis.Client
)

type config struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

func ConnectRedis() {
	var (
		bs     []byte
		err    error
		conf   = &config{}
		logger *zap.Logger
	)

	if bs, err = os.ReadFile("./cache/config.json"); err != nil {
		panic(err)
	}
	if err = jsoniter.Unmarshal(bs, conf); err != nil {
		panic(err)
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password, // no password set
		DB:       0,             // use default DB
	})

	if logger, err = zap.NewProduction(); err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infof("Connect redis URL: %s", conf.Addr)
}
