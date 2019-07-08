package cache

import (
	"time"

	"gopkg.in/redis.v5"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) Get(key string) (string, error) {
	strCmd := r.client.Get(key)
	return strCmd.Val(), strCmd.Err()
}

func (r *RedisCache) Set(key, value string) error {
	statusCmd := r.client.Set(key, value, 0)
	return statusCmd.Err()
}

func (r *RedisCache) SetExpire(key, value string, ttl time.Duration) error {
	statusCmd := r.client.Set(key, value, ttl)
	return statusCmd.Err()
}
func (r *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	ok, err := r.client.SetNX(key, value, expiration).Result()
	return ok, err
}

func (r *RedisCache) SAdd(key string, value interface{}) error {
	strCmd := r.client.SAdd(key,value)
	return strCmd.Err()
}

func (r *RedisCache) SMembers(key string) ([]string,error)  {
	strCmd := r.client.SMembers(key)
	return strCmd.Result()
}

