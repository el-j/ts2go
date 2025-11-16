package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client
type Client struct {
	*redis.Client
}

// Connect creates a new Redis client connection
func Connect(url, password string, db int) (*Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Override password and DB if provided
	if password != "" {
		opts.Password = password
	}
	if db >= 0 {
		opts.DB = db
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{client}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.Client.Close()
}

// Health checks Redis health
func (c *Client) Health(ctx context.Context) error {
	return c.Ping(ctx).Err()
}

// SetSession stores a session with expiration
func (c *Client) SetSession(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Set(ctx, key, value, expiration).Err()
}

// GetSession retrieves a session
func (c *Client) GetSession(ctx context.Context, key string) (string, error) {
	val, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("session not found")
	}
	return val, err
}

// DeleteSession deletes a session
func (c *Client) DeleteSession(ctx context.Context, key string) error {
	return c.Del(ctx, key).Err()
}

// ExtendSession extends session expiration
func (c *Client) ExtendSession(ctx context.Context, key string, expiration time.Duration) error {
	return c.Expire(ctx, key, expiration).Err()
}

// SessionExists checks if a session exists
func (c *Client) SessionExists(ctx context.Context, key string) (bool, error) {
	val, err := c.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// ListSessions lists all sessions matching a pattern
func (c *Client) ListSessions(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	iter := c.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

// IncrementCounter increments a counter
func (c *Client) IncrementCounter(ctx context.Context, key string) (int64, error) {
	return c.Incr(ctx, key).Result()
}

// GetCounter gets a counter value
func (c *Client) GetCounter(ctx context.Context, key string) (int64, error) {
	val, err := c.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// SetWithExpiry sets a key with expiration
func (c *Client) SetWithExpiry(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Set(ctx, key, value, expiration).Err()
}

// ZAdd adds a member to a sorted set
func (c *Client) ZAdd(ctx context.Context, key string, member string, score float64) *redis.IntCmd {
	return c.Client.ZAdd(ctx, key, redis.Z{Score: score, Member: member})
}

// ZRange gets members from sorted set by range
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return c.Client.ZRange(ctx, key, start, stop)
}

// ZRem removes a member from sorted set
func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.Client.ZRem(ctx, key, members...)
}

// ZCard gets the number of members in a sorted set
func (c *Client) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return c.Client.ZCard(ctx, key)
}

// SAdd adds members to a set
func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return c.Client.SAdd(ctx, key, members...)
}

// SMembers gets all members of a set
func (c *Client) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return c.Client.SMembers(ctx, key)
}
