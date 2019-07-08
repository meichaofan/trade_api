package redis

import (
	"fmt"
	"strconv"
	"sync"

	"gopkg.in/redis.v5"

	"trade_api/src/main/conf"
)

var (
	dbs    = make(map[string]*redis.Client)
	mutext sync.Mutex
)

func GetRedis() *redis.Client {
	return cache("base")
}

func cache(name string) *redis.Client {
	if _, ok := dbs[name]; !ok {
		mutext.Lock()
		defer mutext.Unlock()
		config := conf.Conf().Resource.Redis[name]
		var db, _ = strconv.Atoi(config.Db)

		dbs[name] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
			Password: config.Pwd, // no password set
			DB:       db,         // use default DB
		})
	}
	return dbs[name]
}
