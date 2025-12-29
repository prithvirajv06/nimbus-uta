package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/prithvirajv06/nimbus-uta/go/notification/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	return r.client.SetNX(ctx, key, data, expiration).Result()
}

func (r *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Cache patterns
func (r *RedisClient) GetOrSet(ctx context.Context, key string, dest interface{}, expiration time.Duration, fetchFunc func() (interface{}, error)) error {
	// Try to get from cache
	err := r.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	// If not in cache, fetch from source
	if err == redis.Nil {
		data, err := fetchFunc()
		if err != nil {
			return err
		}

		// Store in cache
		if err := r.Set(ctx, key, data, expiration); err != nil {
			return err
		}

		// Marshal and unmarshal to populate dest
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonData, dest)
	}

	return err
}
