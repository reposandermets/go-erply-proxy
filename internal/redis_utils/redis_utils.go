package redis_utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisUtil interface {
	GenerateUniqueKey(r *http.Request) (string, string)
	GetFromCache(ctx context.Context, key string) (string, error)
	SaveToCache(ctx context.Context, key1, key2, data string) error
	ClearCache(ctx context.Context, categoryKey string) error
	FlushRedis(ctx context.Context) error
	PeriodicallyClearCache()
	ManageClearCache(wg *sync.WaitGroup, r *http.Request)
	ManageSaveToCache(wg *sync.WaitGroup, r *http.Request, categoryKey string, urlParamsKey string, jsonData []byte)
}

type RedisUtilImpl struct {
	client *redis.Client
}

func NewRedisUtil(redisClient *redis.Client) RedisUtil {
	return &RedisUtilImpl{
		client: redisClient,
	}
}

func (ru *RedisUtilImpl) GenerateUniqueKey(r *http.Request) (string, string) {
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

func (r *RedisUtilImpl) GetFromCache(ctx context.Context, key string) (string, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return data, nil
}

func (r *RedisUtilImpl) SaveToCache(ctx context.Context, key1, key2, data string) error {
	// Create a composite key
	compositeKey := key1 + "|" + key2
	err := r.client.Set(ctx, compositeKey, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisUtilImpl) ClearCache(ctx context.Context, categoryKey string) error {
	keys, err := r.client.Keys(ctx, categoryKey+"|*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		err = r.client.Del(ctx, keys...).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisUtilImpl) FlushRedis(ctx context.Context) error {
	statusCmd := r.client.FlushAll(ctx)
	return statusCmd.Err()
}

func (r *RedisUtilImpl) PeriodicallyClearCache() {
	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		err := r.FlushRedis(context.Background())
		if err != nil {
			fmt.Println("Error occurred while flushing cache:", err)
		} else {
			fmt.Println("Cache flushed successfully")
		}
	}
}

func (ru *RedisUtilImpl) ManageClearCache(wg *sync.WaitGroup, r *http.Request) {
	defer wg.Done()

	categoryKey, _ := ru.GenerateUniqueKey(r)

	err := ru.ClearCache(r.Context(), categoryKey)
	if err != nil {
		log.Printf("Error clearing cache: %v\n", err)
	}
}

func (ru *RedisUtilImpl) ManageSaveToCache(wg *sync.WaitGroup, r *http.Request, categoryKey string, urlParamsKey string, jsonData []byte) {
	defer wg.Done()

	err := ru.SaveToCache(r.Context(), categoryKey, urlParamsKey, string(jsonData))
	if err != nil {
		log.Printf("Error saving to cache: %v\n", err)
	}
}
