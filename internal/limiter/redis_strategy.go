package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisStrategy struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStrategy(addr string) *RedisStrategy {
	return &RedisStrategy{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		ctx: context.Background(),
	}
}

func (r *RedisStrategy) AllowRequest(key string, limit int, windowSec int) (bool, error) {
	blocked, _ := r.BlockDurationExceeded(key)
	if blocked {
		return false, nil
	}

	now := time.Now().Unix()

	member := uuid.New().String()

	pipe := r.client.TxPipeline()
	pipe.ZAdd(r.ctx, key, redis.Z{Score: float64(now), Member: member})
	pipe.ZRemRangeByScore(r.ctx, key, "0", fmt.Sprintf("%d", now-int64(windowSec)))
	card := pipe.ZCard(r.ctx, key)
	pipe.Expire(r.ctx, key, time.Duration(windowSec)*time.Second)
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return false, err
	}

	count := card.Val()
	return int(count) <= limit, nil
}

func (r *RedisStrategy) BlockDurationExceeded(key string) (bool, error) {
	val, err := r.client.Get(r.ctx, key+":block").Result()
	if err == redis.Nil {
		return false, nil
	}
	return val == "1", err
}

func (r *RedisStrategy) SetBlock(key string, durationSec int) error {
	return r.client.Set(r.ctx, key+":block", "1", time.Duration(durationSec)*time.Second).Err()
}

func (r *RedisStrategy) FlushDB() error {
	return r.client.FlushDB(r.ctx).Err()
}