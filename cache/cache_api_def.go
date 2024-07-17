package cache

import (
	"context"
	"encoding"
	"time"
)

type Reader interface {
	Get(ctx context.Context, k string) (string, error)
	Gets(ctx context.Context, ks []string) (map[string]string, error)
}

type Writer interface {
	Set(ctx context.Context, k string, v encoding.BinaryMarshaler, exp time.Duration) error
	Sets(ctx context.Context, ks []string, vs []encoding.BinaryMarshaler, exps []time.Duration) error
	Deleter
}

type ReadWriter interface {
	Reader
	Writer
}

type Locker interface {
	Lock(ctx context.Context, k string, exp time.Duration) error
	Locks(ctx context.Context, ks []string, exp time.Duration) error
	Deleter
}

type Deleter interface {
	Del(ctx context.Context, k string) error
	Deletes(ctx context.Context, ks []string) error
}
