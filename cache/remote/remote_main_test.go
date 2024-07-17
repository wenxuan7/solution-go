package remote

import (
	"github.com/wenxuan7/solution/link"
	"log/slog"
	"os"
	"testing"
)

var s = NewService()

func setup() {
	link.Redis()
	s.WithWrapper(false, true)
}

func tearDown() {
	err := link.RedisDb.Close()
	if err != nil {
		slog.Error("redis close failed", "error", err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}
