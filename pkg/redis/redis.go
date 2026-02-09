package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func NewRedisClient(addr, password string, db int) *Client {
	var opts *redis.Options
	var err error

	// Try to parse addr as a URL (e.g., redis:// or rediss://)
	opts, err = redis.ParseURL(addr)
	if err != nil {
		// If not a valid URL, use the provided parameters
		opts = &redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}
	}

	rdb := redis.NewClient(opts)
	return &Client{client: rdb}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Client) IncrementWithExpiry(ctx context.Context, key string, expiry time.Duration) (int64, error) {
	pipe := c.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, expiry)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Result()
}
