package cache

import (
	"context"
	"time"
)

type Reader interface {
	Get(ctx context.Context, k string) (string, error)
	Gets(ctx context.Context, ks []string) (map[string]string, error)
}

type Writer interface {
	Set(ctx context.Context, k string, v string, exp time.Duration) error
	Sets(ctx context.Context, ks []string, vs []string, exps []time.Duration) error
	Deleter
}

type Deleter interface {
	Del(ctx context.Context, k string) error
	Deletes(ctx context.Context, ks []string) error
}

type ReadWriter interface {
	Reader
	Writer
}

type Locker interface {
	Lock(ctx context.Context, k string, exp time.Duration) (string, error)
	Locks(ctx context.Context, ks []string, exp time.Duration) ([]string, error)
	UnLock(ctx context.Context, k, v string)
	UnLocks(ctx context.Context, ks, vs []string)
}
