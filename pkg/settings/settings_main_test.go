package settings

import (
	"github.com/wenxuan7/solution/pkg/external"
	"os"
	"testing"
)

var s *Service

func setup() {
	external.Mysql()
	external.Redis()

	var err error
	s, err = NewWithLCache()
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
