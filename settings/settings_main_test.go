package settings

import (
	"github.com/wenxuan7/solution/external"
	"os"
	"testing"
)

var s *Service

func setup() {
	external.Mysql()
	external.Redis()

	var err error
	s, err = NewServiceWithLCache()
	if err != nil {
		panic(err)
	}
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	code := m.Run()
	os.Exit(code)
}
