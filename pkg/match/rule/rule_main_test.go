package rule

import (
	"github.com/wenxuan7/solution/pkg/external"
	"os"
	"testing"
)

func setup() {
	external.Mysql()
	external.Redis()
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	os.Exit(m.Run())
}
