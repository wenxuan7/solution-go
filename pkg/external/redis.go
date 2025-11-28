package external

type redisConf struct {
	Addr     string
	Password string
}

func newRedisConf() *redisConf {
	return &redisConf{
		Addr:     "localhost:6379",
		Password: "wenxuan101314",
	}
}
