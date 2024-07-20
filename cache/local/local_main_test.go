package local

import (
	"github.com/allegro/bigcache"
	"github.com/wenxuan7/solution/external"
	"os"
	"testing"
	"time"
)

var s *Service

func setup() {
	external.Redis()
	var err error
	s, err = NewService(true, "bigCacheChannel", bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	s.SubscribeClose()
	err := s.lCache.Close()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	err = external.RedisDb.Close()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
