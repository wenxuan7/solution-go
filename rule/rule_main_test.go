package rule

import (
	"os"
	"testing"
)

func setup() {

}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	os.Exit(m.Run())
}
