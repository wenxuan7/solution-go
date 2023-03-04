package cache

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/solution-go/log"
	"os"
	"path"
	"runtime"
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
		bs      []byte
		currDir string
		err     error
		conf    = &config{}
	)

	_, currFile, _, _ := runtime.Caller(0)
	currDir = path.Dir(currFile)
	if bs, err = os.ReadFile(currDir + "/config.json"); err != nil {
		log.Sugar.Panic(err)
	}
	if err = jsoniter.Unmarshal(bs, conf); err != nil {
		log.Sugar.Panic()
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password, // no password set
		DB:       0,             // use default DB
	})

	log.Sugar.Infof("Connect redis URL: %s", conf.Addr)
}
