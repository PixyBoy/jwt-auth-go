package redisadp

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(addr, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		MinIdleConns: 4,
		PoolSize:     16,
	})
	// اتصال تست سبک
	_ = rdb.Ping(context.Background()).Err()
	return rdb
}
