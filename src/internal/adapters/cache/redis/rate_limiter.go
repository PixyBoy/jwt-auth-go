package redisadp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiterRedis struct {
	rdb    *redis.Client
	prefix string
}

func NewRateLimiter(rdb *redis.Client, prefix string) *RateLimiterRedis {
	if prefix == "" {
		prefix = "auth:otp:rl:"
	}
	return &RateLimiterRedis{rdb: rdb, prefix: prefix}
}

func (l *RateLimiterRedis) key(key string) string {
	return fmt.Sprintf("%s%s", l.prefix, key)
}

func (l *RateLimiterRedis) Allow(key string, max int, windowSeconds int) (bool, int, error) {
	ctx := context.Background()
	rkey := l.key(key)

	// INCR
	cnt, err := l.rdb.Incr(ctx, rkey).Result()
	if err != nil {
		return false, 0, err
	}
	// First TTL
	if cnt == 1 {
		_ = l.rdb.Expire(ctx, rkey, time.Duration(windowSeconds)*time.Second).Err()
	}
	allowed := cnt <= int64(max)
	return allowed, int(cnt), nil
}
