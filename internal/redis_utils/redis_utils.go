package redis_utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// GenerateUniqueKey generates a unique key from the request URL.
func GenerateUniqueKey(r *http.Request) (string, string) {
	// Get the URI path
	uriPath := r.URL.Path

	// Get the sorted URL parameters
	params := r.URL.Query()
	paramKeys := make([]string, 0, len(params))
	for key := range params {
		paramKeys = append(paramKeys, key)
	}
	sort.Strings(paramKeys)

	var paramList []string
	for _, key := range paramKeys {
		values := params[key]
		for _, value := range values {
			paramList = append(paramList, fmt.Sprintf("%s=%s", key, value))
		}
	}

	key1 := uriPath
	key2 := strings.Join(paramList, "_")
	return key1, key2
}

// GetFromCache retrieves data from Redis cache using the unique key.
func GetFromCache(ctx context.Context, key string) (string, error) {
	// Retrieve the Redis client from the context
	client := ctx.Value("redis").(*redis.Client)

	// Check if the key exists in the cache
	data, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Key does not exist in the cache
			return "", nil
		}
		// Error occurred while accessing Redis
		return "", err
	}

	// Key exists in the cache, return the data
	return data, nil
}

// SaveToCache saves data to Redis cache using the composite keys.
func SaveToCache(ctx context.Context, key1, key2, data string) error {
	// Retrieve the Redis client from the context
	client := ctx.Value("redis").(*redis.Client)

	// Create a composite key
	compositeKey := key1 + "|" + key2

	// Save data to Redis cache with the composite key
	err := client.Set(ctx, compositeKey, data, 0).Err()
	if err != nil {
		// Error occurred while saving to Redis
		return err
	}

	// Data successfully saved to cache
	return nil
}

// GetRedisClientFromContext retrieves the Redis client from the request context.
func GetRedisClientFromContext(ctx context.Context) (*redis.Client, error) {
	client, ok := ctx.Value("redis").(*redis.Client)
	if !ok {
		return nil, errors.New("redis client not found in context")
	}
	return client, nil
}

// ClearCache clears the cache by the key.
func ClearCache(ctx context.Context, client *redis.Client, key string) error {
	// Find keys matching the provided key pattern
	keys, err := client.Keys(ctx, key+"|*").Result()
	if err != nil {
		// Error occurred while accessing Redis
		return err
	}

	// Delete the keys matching the pattern
	if len(keys) > 0 {
		err = client.Del(ctx, keys...).Err()
		if err != nil {
			// Error occurred while deleting keys from Redis
			return err
		}
	}

	// Cache cleared successfully
	return nil
}

// FlushRedis flushes (empties) the entire Redis database.
func FlushRedis(ctx context.Context, client *redis.Client) error {
	statusCmd := client.FlushAll(ctx)
	return statusCmd.Err()
}

// Periodically clear the cache every 10 seconds
func PeriodicallyClearCache(client *redis.Client) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		// Clear cache by the key ("/v1/brand")
		err := FlushRedis(context.Background(), client)
		if err != nil {
			fmt.Println("Error occurred while flushing cache:", err)
		} else {
			fmt.Println("Cache flushed successfully")
		}
	}
}

func WithRedisContext(handler http.Handler, client *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new context with the Redis client
		ctx := context.WithValue(r.Context(), "redis", client)

		// Serve the request with the new context
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
