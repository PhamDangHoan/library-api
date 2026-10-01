package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}

type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache service.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Get retrieves a cached value and unmarshals it into dest.
func (r *RedisCache) Get(
	ctx context.Context,
	key string,
	dest interface{},
) error {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal(
		[]byte(value),
		dest,
	)
}

// Set stores a value in Redis with a TTL.
func (r *RedisCache) Set(
	ctx context.Context,
	key string,
	value interface{},
	ttl time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(
		ctx,
		key,
		data,
		ttl,
	).Err()
}

// Delete deletes one cache key.
func (r *RedisCache) Delete(
	ctx context.Context,
	key string,
) error {
	return r.client.Del(
		ctx,
		key,
	).Err()
}

// DeleteByPattern deletes all keys matching a pattern.
//
// Example:
//
//	"library:books:*"
//
// This can be used to invalidate all book list caches
// after creating, updating or deleting a book.
func (r *RedisCache) DeleteByPattern(
	ctx context.Context,
	pattern string,
) error {
	var cursor uint64

	for {
		keys, nextCursor, err := r.client.Scan(
			ctx,
			cursor,
			pattern,
			100,
		).Result()

		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := r.client.Del(
				ctx,
				keys...,
			).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor

		if cursor == 0 {
			break
		}
	}

	return nil
}
