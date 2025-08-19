package redisadp

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPStoreRedis struct {
	rdb    *redis.Client
	prefix string
}

func NewOTPStore(rdb *redis.Client, prefix string) *OTPStoreRedis {
	if prefix == "" {
		prefix = "auth:otp:"
	}
	return &OTPStoreRedis{rdb: rdb, prefix: prefix}
}

func (s *OTPStoreRedis) key(phone string) string {
	return s.prefix + phone
}

// Save: hash, attempts=0, issued_at=now; EXPIRE ttlSeconds
func (s *OTPStoreRedis) Save(phone, hash string, ttlSeconds int) error {
	ctx := context.Background()
	key := s.key(phone)
	pipe := s.rdb.TxPipeline()
	pipe.HSet(ctx, key, map[string]interface{}{
		"hash":      hash,
		"attempts":  0,
		"issued_at": time.Now().Unix(),
	})
	pipe.Expire(ctx, key, time.Duration(ttlSeconds)*time.Second)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *OTPStoreRedis) Get(phone string) (hash string, attempts int, exists bool, err error) {
	ctx := context.Background()
	key := s.key(phone)

	res, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return "", 0, false, err
	}
	if len(res) == 0 {
		return "", 0, false, nil
	}

	hash = res["hash"]
	if aStr, ok := res["attempts"]; ok {
		if v, convErr := strconv.Atoi(aStr); convErr == nil {
			attempts = v
		}
	}
	return hash, attempts, true, nil
}

func (s *OTPStoreRedis) IncreaseAttempt(phone string) (int, error) {
	ctx := context.Background()
	key := s.key(phone)
	val, err := s.rdb.HIncrBy(ctx, key, "attempts", 1).Result()
	return int(val), err
}

func (s *OTPStoreRedis) Delete(phone string) error {
	ctx := context.Background()
	key := s.key(phone)
	return s.rdb.Del(ctx, key).Err()
}
