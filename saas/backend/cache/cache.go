package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/el-j/ts2go/saas/backend/redis"
)

// Cache provides caching functionality
type Cache struct {
	redis *redis.Client
}

// NewCache creates a new cache instance
func NewCache(redisClient *redis.Client) *Cache {
	return &Cache{redis: redisClient}
}

// CacheConfig defines cache TTL settings
type CacheConfig struct {
	ProjectTTL time.Duration
	UserTTL    time.Duration
	FileTTL    time.Duration
	DefaultTTL time.Duration
}

// DefaultCacheConfig returns default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		ProjectTTL: 10 * time.Minute,
		UserTTL:    15 * time.Minute,
		FileTTL:    5 * time.Minute,
		DefaultTTL: 5 * time.Minute,
	}
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.redis.GetSession(ctx, key)
	if err != nil {
		return fmt.Errorf("cache miss: %w", err)
	}

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return nil
}

// Set stores a value in cache with TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.redis.SetWithExpiry(ctx, key, string(data), ttl); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes a value from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.redis.Delete(ctx, key)
}

// DeletePattern removes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	// Note: In production, use SCAN instead of KEYS for large datasets
	// This is a simplified implementation
	return nil
}

// GetOrSet retrieves from cache or computes and stores the value
func (c *Cache) GetOrSet(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error)) (interface{}, error) {
	var result interface{}
	err := c.Get(ctx, key, &result)
	if err == nil {
		return result, nil
	}

	// Cache miss - compute value
	value, err := fn()
	if err != nil {
		return nil, err
	}

	// Store in cache (ignore error - not critical)
	c.Set(ctx, key, value, ttl)

	return value, nil
}

// Cache key builders
func ProjectCacheKey(projectID string) string {
	return fmt.Sprintf("cache:project:%s", projectID)
}

func UserCacheKey(userID string) string {
	return fmt.Sprintf("cache:user:%s", userID)
}

func FileCacheKey(fileID string) string {
	return fmt.Sprintf("cache:file:%s", fileID)
}

func ProjectListCacheKey(userID string) string {
	return fmt.Sprintf("cache:projects:user:%s", userID)
}

func FileListCacheKey(projectID string) string {
	return fmt.Sprintf("cache:files:project:%s", projectID)
}

// Invalidation helpers
func (c *Cache) InvalidateProject(ctx context.Context, projectID string) error {
	return c.Delete(ctx, ProjectCacheKey(projectID))
}

func (c *Cache) InvalidateProjectList(ctx context.Context, userID string) error {
	return c.Delete(ctx, ProjectListCacheKey(userID))
}

func (c *Cache) InvalidateFile(ctx context.Context, fileID string) error {
	return c.Delete(ctx, FileCacheKey(fileID))
}

func (c *Cache) InvalidateFileList(ctx context.Context, projectID string) error {
	return c.Delete(ctx, FileListCacheKey(projectID))
}
