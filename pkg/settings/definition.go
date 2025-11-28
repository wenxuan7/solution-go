package settings

import (
	"context"
	"fmt"
	"github.com/wenxuan7/solution/pkg/utils"
)

type Reader interface {
	Get(ctx context.Context, key string) (string, error)
	Gets(ctx context.Context, keys []string) (map[string]string, error)
	GetFromDb(ctx context.Context, key string) (*Settings, error)
	GetsFromDb(ctx context.Context, keys []string) (map[string]*Settings, error)
	GetKV(ctx context.Context, k string) (KV, error)
	GetsKV(ctx context.Context, ks []string) (map[string]KV, error)
}

type Writer interface {
	Set(ctx context.Context, e *Settings) error
	Sets(ctx context.Context, es []*Settings) error
}

type ReadWriter interface {
	Reader
	Writer
}

type Settings struct {
	utils.Model
	CompanyId  uint   `json:"company_id"`
	K          string `json:"k"`
	V          string `json:"v"`
	CreatedMsg string `json:"created_msg"`
}

func (e *Settings) TableName() string {
	return "settings"
}

type Context struct {
	ks map[string]struct{}
	vs map[string]KV
}

func NewContext(ks ...string) *Context {
	return nil
}

func (c *Context) Get(k string) (KV, error) {
	entity, ok := c.ks[k]
	if !ok {
		return nil, fmt.Errorf("settings: no such key: %s", k)
	}
	return entity, nil
}
