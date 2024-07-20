package settings

import (
	"context"
	"github.com/wenxuan7/solution/utils"
)

type Reader interface {
	Get(ctx context.Context, key string) (string, error)
	Gets(ctx context.Context, keys []string) (map[string]string, error)
	GetFromDb(ctx context.Context, key string) (*Entity, error)
	GetsFromDb(ctx context.Context, keys []string) (map[string]*Entity, error)
	GetKV(ctx context.Context, k string) (KV, error)
	GetsKV(ctx context.Context, ks []string) (map[string]KV, error)
}

type Writer interface {
	Set(ctx context.Context, e *Entity) error
	Sets(ctx context.Context, es []*Entity) error
}

type ReadWriter interface {
	Reader
	Writer
}

type Entity struct {
	utils.Model
	CompanyId  uint   `json:"company_id"`
	K          string `json:"k"`
	V          string `json:"v"`
	CreatedMsg string `json:"created_msg"`
}

func (e *Entity) TableName() string {
	return "settings"
}
