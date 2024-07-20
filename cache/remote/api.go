package remote

import (
	"context"
	"fmt"
	"github.com/wenxuan7/solution/external"
	"github.com/wenxuan7/solution/utils"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// Service 提供redis操作api
type Service struct {
	wrapperEnv       bool
	wrapperCompanyId bool
}

// NewService Service构造器
func NewService() *Service {
	return &Service{}
}

func (s *Service) WithWrapper(wrapperEnv, wrapperCompany bool) *Service {
	s.wrapperEnv = wrapperEnv
	s.wrapperCompanyId = wrapperCompany
	return s
}

// wrapper 包装k
func (s *Service) wrapper(ctx context.Context, k string) string {
	sb := strings.Builder{}
	if s.wrapperCompanyId {
		sb.WriteString(strconv.FormatUint(uint64(utils.GetCompanyId(ctx)), 10))
		sb.WriteString("_")
	}
	if s.wrapperEnv {
		sb.WriteString(ctx.Value("env").(string))
		sb.WriteString("_")
	}
	sb.WriteString(k)
	return sb.String()
}

// Set 单个set
func (s *Service) Set(ctx context.Context, k string, v string, exp time.Duration) error {
	if exp <= 0 {
		return fmt.Errorf("remote: invalid expiration")
	}
	wrapperK := s.wrapper(ctx, k)
	_, err := external.RedisDb.Set(ctx, wrapperK, v, exp).Result()
	if err != nil {
		return fmt.Errorf("remote: failed to set value for wrapperK '%s': %w", wrapperK, err)
	}
	return nil
}

var luaSets = `
    local setKeys = {}  -- 存储已成功设置的键

    for i = 1, #KEYS do
        local key = KEYS[i]
        local value = ARGV[(i-1) * 2 + 1]
        local ttl = tonumber(ARGV[(i-1) * 2 + 2])

        -- 尝试设置键的值和过期时间
        local result = redis.call("SET", key, value)
        if result.ok ~= "OK" then
            -- 回滚已成功设置的键
            for j = 1, #setKeys do
                redis.call("DEL", setKeys[j])
            end
            return redis.error_reply("Failed to set key " .. key)
        end

        -- 设置过期时间
        result = redis.call("EXPIRE", key, ttl)
        if result == 0 then
            -- 回滚已成功设置的键
            for j = 1, #setKeys do
                redis.call("DEL", setKeys[j])
            end
            return redis.error_reply("Failed to set expiration for key " .. key)
        end

        table.insert(setKeys, key)
    end

    return "OK"
    `

// Sets 批量set 失败后回滚已设置的key
func (s *Service) Sets(ctx context.Context, ks []string, vs []string, exps []time.Duration) error {
	l := len(ks)
	if l > 100 || l != len(vs) || l != len(exps) {
		return fmt.Errorf("remote: invalid length of ks")
	}
	for _, exp := range exps {
		if exp <= 0 {
			return fmt.Errorf("remote: invalid exps")
		}
	}
	keys := make([]string, l)
	for i, v := range ks {
		keys[i] = s.wrapper(ctx, v)
	}
	values := make([]string, l)
	for i, v := range vs {
		values[i] = v
	}

	// Prepare arguments for the Lua script
	args := make([]any, 0, l*2)
	for i := range values {
		args = append(args, values[i], exps[i])
	}

	_, err := external.RedisDb.Eval(ctx, luaSets, keys, args...).Result()
	if err != nil {
		return fmt.Errorf("remote: failed to set values for luaScript: %w", err)
	}
	return nil
}

// Get 单个get
func (s *Service) Get(ctx context.Context, k string) (string, error) {
	wrapperK := s.wrapper(ctx, k)
	ret, err := external.RedisDb.Get(ctx, wrapperK).Result()
	if err != nil {
		return "", fmt.Errorf("remote: failed to get value for wrapperK '%s': %w", wrapperK, err)
	}
	return ret, nil
}

var luaGets = `
        local result = {}
        for i, key in ipairs(KEYS) do
            local value = redis.call("GET", key)
            if not value then
                value = "nil"
            end
            result[i] = value
        end
        return result
    `

// Gets 批量get
func (s *Service) Gets(ctx context.Context, ks []string) (map[string]string, error) {
	l := len(ks)
	if l > 100 {
		return nil, fmt.Errorf("remote: invalid length of ks")
	}
	wrapperKs := make([]string, l)
	for i, v := range ks {
		wrapperKs[i] = s.wrapper(ctx, v)
	}

	ret, err := external.RedisDb.Eval(ctx, luaGets, wrapperKs).Result()
	if err != nil {
		return nil, fmt.Errorf("remote: failed to get values for luaScript: %w", err)
	}
	retArr := ret.([]any)
	retMap := make(map[string]string, len(retArr))
	for i, v := range retArr {
		str := v.(string)
		if str == "nil" {
			continue
		}
		retMap[ks[i]] = str
	}
	return retMap, nil
}

// Del 单个删除
func (s *Service) Del(ctx context.Context, k string) error {
	wrapperK := s.wrapper(ctx, k)
	_, err := external.RedisDb.Del(ctx, wrapperK).Result()
	if err != nil {
		return fmt.Errorf("remote: failed to delete value for wrapperK '%s': %w", wrapperK, err)
	}
	return nil
}

// DDel 延迟双删
func (s *Service) DDel(ctx context.Context, k string) error {
	err := s.Del(ctx, k)
	if err != nil {
		return fmt.Errorf("remote: failed to double delete key '%s' on first: %w", k, err)
	}
	go func(k string) {
		time.Sleep(time.Second * 3)
		err2 := s.Del(ctx, k)
		if err2 != nil {
			endErr := fmt.Errorf("remote: failed to double delete key '%s' on second: %v", k, err2)
			slog.Error("remote: DDel", "error", endErr)
		}
	}(k)
	return nil
}

var luaDeletes = `
    local keys = KEYS
    for i = 1, #keys do
        redis.call("DEL", keys[i])
    end
    return #keys
    `

// Deletes 批量删除
func (s *Service) Deletes(ctx context.Context, ks []string) error {
	l := len(ks)
	if l > 100 {
		return fmt.Errorf("remote: invalid length of ks")
	}
	wrapperKs := make([]string, len(ks))
	for i, v := range ks {
		wrapperKs[i] = s.wrapper(ctx, v)
	}
	_, err := external.RedisDb.Eval(ctx, luaDeletes, wrapperKs).Result()
	if err != nil {
		return fmt.Errorf("remote: failed to delete values for luaScript: %w", err)
	}
	return nil
}

// DDeletes 批量延迟删除
func (s *Service) DDeletes(ctx context.Context, ks []string) error {
	err := s.Deletes(ctx, ks)
	if err != nil {
		return fmt.Errorf("remote: failed to double deletes keys '%s' on first: %w", ks, err)
	}
	go func(ks []string) {
		time.Sleep(time.Second * 3)
		err2 := s.Deletes(ctx, ks)
		if err2 != nil {
			endErr := fmt.Errorf("remote: failed to double deletes keys '%s' on second: %v", ks, err2)
			slog.Error("remote: DDeletes", "error", endErr)
		}
	}(ks)
	return nil
}

var luaLock = `
    local key = KEYS[1]
    local ttl = tonumber(ARGV[1])
    local lockValue = ARGV[2]
	-- 尝试设置键的值和过期时间，如果键已存在则返回错误
	local result = redis.call("SET", key, lockValue, "NX", "EX", ttl)
	if result == false then
		return redis.error_reply("Key " .. key .. " is already locked")
	else
		return "OK"
	end
    `

// Lock 单个加锁
func (s *Service) Lock(ctx context.Context, k string, exp time.Duration) error {
	if exp <= 0 {
		return fmt.Errorf("remote: invalid length of exp")
	}
	wrapperK := s.wrapper(ctx, k)
	err := external.RedisDb.Eval(ctx, luaLock, []string{wrapperK}, exp, 1).Err()
	if err != nil {
		return fmt.Errorf("remote: failed to setNx for wrapperK '%s': %w", wrapperK, err)
	}
	return nil
}

var luaLocks = `
    local keys = KEYS
    local ttl = tonumber(ARGV[1])
    local lockValue = ARGV[2]
    local lockedKeys = {}  -- 存储已成功加锁的键

    for i = 1, #keys do
        -- 尝试设置键的值和过期时间，如果键已存在则返回错误
        local result = redis.call("SET", keys[i], lockValue, "NX", "EX", ttl)
        if result == false then
            -- 回滚已成功加锁的键
            for j = 1, #lockedKeys do
                redis.call("DEL", lockedKeys[j])
            end
            return redis.error_reply("Key " .. keys[i] .. " is already locked")
        else
            table.insert(lockedKeys, keys[i])
        end
    end
    return "OK"
    `

// Locks 批量加锁 任意key失败 回滚加锁成功的key
func (s *Service) Locks(ctx context.Context, ks []string, exp time.Duration) error {
	l := len(ks)
	if l > 100 {
		return fmt.Errorf("remote: invalid length of ks")
	}
	if exp <= 0 {
		return fmt.Errorf("remote: invalid length of exp")
	}
	wrapperKs := make([]string, len(ks))
	for i, v := range ks {
		wrapperKs[i] = s.wrapper(ctx, v)
	}

	_, err := external.RedisDb.Eval(ctx, luaLocks, wrapperKs, exp, 1).Result()
	if err != nil {
		return fmt.Errorf("remote: failed to lock values for luaScript: %w", err)
	}
	return nil
}
