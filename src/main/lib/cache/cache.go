package cache

import (
	"time"
)

type Cache interface {
	Get(string) (string, error)
	Set(string, string) error
	SetExpire(string, string, time.Duration) error
	SetNX(key string, value interface{}, expiration time.Duration) (bool, error)
	SAdd(key string, value interface{}) error
	SMembers(key string) ([]string,error)
}
