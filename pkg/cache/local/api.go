package local

import (
	"context"
	"errors"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/redis/go-redis/v9"
	"github.com/wenxuan7/solution/pkg/cache"
	"github.com/wenxuan7/solution/pkg/cache/remote"
	"github.com/wenxuan7/solution/pkg/external"
	"github.com/wenxuan7/solution/pkg/utils"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	remoteRW         cache.ReadWriter
	lCache           *bigcache.BigCache
	wrapperCompanyId bool
	channel          string
	channelCtx       context.Context
	channelCtxCancel context.CancelFunc
}

type keysChannel []string

func (kc *keysChannel) MarshalBinary() (data []byte, err error) {
	sb := strings.Builder{}
	sb.WriteString("[")
	l := len(*kc)
	for i, key := range *kc {
		sb.WriteString(key)
		if i != l-1 {
			sb.WriteString(",")
		} else {
			sb.WriteString("]")
		}
	}
	return []byte(sb.String()), nil
}

func (kc *keysChannel) UnmarshalBinary(data []byte) error {
	input := strings.Trim(string(data), "[]")
	*kc = strings.Split(input, ",")
	return nil
}

func NewService(wrapperCompanyId bool, channel string, config bigcache.Config) (*Service, error) {
	if channel == "" {
		return nil, errors.New("local: channel is required")
	}
	lCache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, fmt.Errorf("local: fail to new bigCache; %w", err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	s := &Service{
		remoteRW:         remote.NewService().WithWrapper(false, true),
		wrapperCompanyId: wrapperCompanyId,
		channel:          channel,
		lCache:           lCache,
		channelCtx:       ctx,
		channelCtxCancel: cancelFunc,
	}

	go s.subscribe()
	return s, nil
}

func (s *Service) wrapper(ctx context.Context, k string) string {
	sb := strings.Builder{}
	if s.wrapperCompanyId {
		sb.WriteString(strconv.FormatUint(uint64(utils.GetCompanyId(ctx)), 10))
		sb.WriteString("_")
	}
	sb.WriteString(k)
	return sb.String()
}

func (s *Service) Get(ctx context.Context, k string) (string, error) {
	wrapperK := s.wrapper(ctx, k)
	bs, err := s.lCache.Get(wrapperK)
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		var str string
		str, err = s.remoteRW.Get(ctx, k)
		if err != nil && !errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("local: fail redis to get for wrapper key '%s'; %w", wrapperK, err)
		}
		err = s.lCache.Set(wrapperK, []byte(str))
		if err != nil {
			return "", fmt.Errorf("local: fail bigCache to set for wrapper key '%s'; %w", wrapperK, err)
		}
		return str, nil
	}
	if err != nil {
		return "", fmt.Errorf("local: fail bigCache to get for key '%s':%w", k, err)
	}
	return string(bs), nil
}

func (s *Service) Gets(ctx context.Context, ks []string) (map[string]string, error) {
	l := len(ks)
	wrapperKs := make([]string, l)
	for i, v := range ks {
		wrapperKs[i] = s.wrapper(ctx, v)
	}
	ret := make(map[string]string, l)
	remoteKs := make([]string, 0, l)
	for i, v := range wrapperKs {
		value, err := s.lCache.Get(v)
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			remoteKs = append(remoteKs, ks[i])
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("local: fail bigCache to get for wrapper key '%s'; %w", wrapperKs[i], err)
		}
		ret[ks[i]] = string(value)
	}
	if len(remoteKs) == 0 {
		return ret, nil
	}
	remoteValues, err := s.remoteRW.Gets(ctx, remoteKs)
	if err != nil {
		return nil, fmt.Errorf("local: fail redis to gets for wrapper keys; %w", err)
	}
	for k, v := range remoteValues {
		ret[k] = v
	}
	return ret, nil
}

func (s *Service) Set(ctx context.Context, k string, v string, exp time.Duration) error {
	err := s.remoteRW.Set(ctx, k, v, exp)
	if err != nil {
		return fmt.Errorf("local: fail remote to set; %w", err)
	}
	err = s.publish(ctx, []string{k})
	if err != nil {
		return fmt.Errorf("local: fail publish keys in set for channel '%s'; %w", s.channel, err)
	}
	return nil
}

func (s *Service) Sets(ctx context.Context, ks []string, vs []string, exps []time.Duration) error {
	err := s.remoteRW.Sets(ctx, ks, vs, exps)
	if err != nil {
		return fmt.Errorf("local: fail remote to sets; %w", err)
	}
	err = s.publish(ctx, ks)
	if err != nil {
		return fmt.Errorf("local: fail publish keys int sets for channel '%s'; %w", s.channel, err)
	}
	return nil
}

func (s *Service) Del(ctx context.Context, k string) error {
	err := s.remoteRW.Del(ctx, k)
	if err != nil {
		return fmt.Errorf("local: fail remote to del; %w", err)
	}
	err = s.publish(ctx, []string{k})
	if err != nil {
		return fmt.Errorf("local: fail publish keys in del for channel '%s'; %w", s.channel, err)
	}
	return nil
}

func (s *Service) Deletes(ctx context.Context, ks []string) error {
	err := s.remoteRW.Deletes(ctx, ks)
	if err != nil {
		return fmt.Errorf("local: fail remote to deletes; %w", err)
	}
	err = s.publish(ctx, ks)
	if err != nil {
		return fmt.Errorf("local: fail publish keys in deletes for channel '%s'; %w", s.channel, err)
	}
	return nil
}

var ErrPublishRedis = errors.New("local: fail publish keys")

func (s *Service) publish(ctx context.Context, ks []string) error {
	wrapperKs := make(keysChannel, 0, len(ks))
	for _, v := range ks {
		wrapperKs = append(wrapperKs, s.wrapper(ctx, v))
	}
	_, err := external.RedisDb.Publish(ctx, s.channel, &wrapperKs).Result()
	if err != nil {
		return fmt.Errorf("local: wrapper keys '%s', channel '%s'; %w; %w", wrapperKs, s.channel, ErrPublishRedis, err)
	}
	return nil
}

func (s *Service) subscribe() {
	sub := external.RedisDb.Subscribe(context.Background(), s.channel)
	defer func() {
		err := sub.Close()
		if err != nil {
			slog.Error("local: fail close redis subscription", "error", err, "channel", s.channel)
		} else {
			slog.Info("local: redis subscription is closed", "channel", s.channel)
		}
	}()

	ch := sub.Channel()

	for {
		select {
		case <-s.channelCtx.Done():
			slog.Info("local: subscription closed by cancelCtx", "subscribeChannel", s.channel)
			return
		case msg, ok := <-ch:
			if !ok {
				slog.Info("local: subscription closed by chan", "subscribeChannel", s.channel)
				return
			}
			slog.Info("local: received message: ", "msg.Payload", msg.Payload)
			var ksChannel keysChannel
			_ = ksChannel.UnmarshalBinary([]byte(msg.Payload))
			for _, k := range ksChannel {
				err := s.lCache.Delete(k)
				if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
					slog.Error("local: fail bigCache to delete", "key", k, "error", err)
				}
			}
		}
	}
}

func (s *Service) SubscribeClose() {
	s.channelCtxCancel()
}
