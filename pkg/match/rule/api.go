package rule

import (
	"github.com/wenxuan7/solution/pkg/cache"
	"github.com/wenxuan7/solution/pkg/cache/remote"
)

type Service struct {
	cache cache.ReadWriter
}

func NewServiceWithRemoteCache() *Service {
	remoteCache := remote.NewService().WithWrapper(false, true)
	return &Service{remoteCache}
}
