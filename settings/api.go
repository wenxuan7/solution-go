package settings

import (
	"context"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/wenxuan7/solution/cache"
	"github.com/wenxuan7/solution/cache/local"
	"github.com/wenxuan7/solution/external"
	"github.com/wenxuan7/solution/utils"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type Service struct {
	cache cache.ReadWriter
}

func NewServiceWithLCache() (*Service, error) {
	lCache, err := local.NewService(true, "settingsLocalCacheChannel", bigcache.DefaultConfig(time.Hour*2))
	if err != nil {
		return nil, fmt.Errorf("settings: fail to local.NewService in NewServiceWithLCache: %w", err)
	}
	return &Service{cache: lCache}, nil
}

func (s *Service) GetKV(ctx context.Context, k string) (KV, error) {
	v, err := s.Get(ctx, k)
	if err != nil {
		return nil, fmt.Errorf("settings: fail to Get in GetKV for k '%s': %w", k, err)
	}
	kv, err := StrconvKV(k, v)
	if err != nil {
		return nil, fmt.Errorf("settings: fail to StrconvKV in GetKV for k '%s': %w", k, err)
	}
	return kv, nil
}

func (s *Service) GetsKV(ctx context.Context, ks []string) (map[string]KV, error) {
	vs, err := s.Gets(ctx, ks)
	if err != nil {
		return nil, fmt.Errorf("settings: fail to Gets in GetsKV for keys '%s': %w", ks, err)
	}
	kvs := make(map[string]KV, len(vs))
	for k, v := range vs {
		kv, err := StrconvKV(k, v)
		if err != nil {
			return nil, fmt.Errorf("settings: fail to StrconvKV in GetsKV for fail k '%s': %w", k, err)
		}
		kvs[k] = kv
	}
	return kvs, nil
}

func (s *Service) Get(ctx context.Context, key string) (string, error) {
	cacheRes, err := s.cache.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("settings: fail to cache.Get in Get for key '%s': %w", key, err)
	}
	if cacheRes != "" {
		return cacheRes, nil
	}

	e, err := s.GetFromDb(ctx, key)
	if err != nil {
		return "", fmt.Errorf("settings: fail to GetFromDb in Get for key '%s': %w", key, err)
	}
	err = s.cache.Set(ctx, e.K, e.V, randExp())
	if err != nil {
		return "", fmt.Errorf("settings: fail to cache.Set in Get for key '%s': %w", key, err)
	}
	return e.V, nil
}

func (s *Service) Gets(ctx context.Context, keys []string) (map[string]string, error) {
	res := make(map[string]string, len(keys))
	cacheRes, err := s.cache.Gets(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("settings: fail to cache.Gets in Gets: %w", err)
	}

	dbKeys := make([]string, 0, len(cacheRes))
	for _, k := range keys {
		v, ok := cacheRes[k]
		if !ok {
			dbKeys = append(dbKeys, k)
			continue
		}
		res[k] = v
	}

	dbRes, err := s.GetsFromDb(ctx, dbKeys)
	if err != nil {
		return nil, fmt.Errorf("settings: fail to GetFromDb in Gets: %w", err)
	}

	newKs := make([]string, 0, len(dbRes))
	newVs := make([]string, 0, len(dbRes))
	newExps := make([]time.Duration, 0, len(dbRes))
	for k, v := range dbRes {
		newKs = append(newKs, k)
		newVs = append(newVs, v.V)
		newExps = append(newExps, randExp())
		res[k] = v.V
	}
	err = s.cache.Sets(ctx, newKs, newVs, newExps)
	if err != nil {
		return nil, fmt.Errorf("settings: fail to cache.Sets in Gets: %w", err)
	}

	return res, nil
}

func (s *Service) GetFromDb(ctx context.Context, key string) (*Entity, error) {
	companyId := utils.GetCompanyId(ctx)
	e := &Entity{K: key, CompanyId: companyId}
	err := external.MysqlDb.Model(e).First(e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return e, nil
	}
	if err != nil {
		return nil, fmt.Errorf("settings: fail to query first in GetFromDb for key '%s': %w", key, err)
	}
	return e, nil
}

func (s *Service) GetsFromDb(ctx context.Context, keys []string) (map[string]*Entity, error) {
	companyId := utils.GetCompanyId(ctx)
	es := make([]*Entity, len(keys))
	res := make(map[string]*Entity, len(keys))
	err := external.MysqlDb.Where("k IN ?", keys).Find(&es).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		for _, k := range keys {
			res[k] = &Entity{K: k, CompanyId: companyId}
		}
		return res, nil
	}
	if err != nil {
		return nil, fmt.Errorf("settings: fail to Find in GetsFromDb for keys '%s': %w", keys, err)
	}

	for _, e := range es {
		res[e.K] = e
	}
	for _, k := range keys {
		if _, ok := res[k]; !ok {
			res[k] = &Entity{K: k, CompanyId: companyId}
		}
	}
	return res, nil
}

func (s *Service) Set(ctx context.Context, e *Entity) error {
	companyId := utils.GetCompanyId(ctx)
	e.CompanyId = companyId
	err := external.MysqlDb.Save(e).Error
	if err != nil {
		return fmt.Errorf("settings: fail to Save in Set: %w", err)
	}

	err = s.cache.Del(ctx, e.K)
	if err != nil {
		return fmt.Errorf("settings: fail to s.cache.Del in Set: %w", err)
	}
	return nil
}

func (s *Service) Sets(ctx context.Context, es []*Entity) error {
	companyId := utils.GetCompanyId(ctx)
	ks := make([]string, 0, len(es))
	for _, e := range es {
		e.CompanyId = companyId
		ks = append(ks, e.K)
	}
	err := external.MysqlDb.Save(es).Error
	if err != nil {
		return fmt.Errorf("settings: fail to Save in Sets: %w", err)
	}

	err = s.cache.Deletes(ctx, ks)
	if err != nil {
		return fmt.Errorf("settings: fail to s.cache.Deletes in Sets: %w", err)
	}
	return nil
}

// randExp 默认2小时 + 随机秒数1000
func randExp() time.Duration {
	n := rand.Intn(1000)
	return time.Hour*2 + time.Second*time.Duration(n)
}
