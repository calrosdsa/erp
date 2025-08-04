package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v9"

	"github.com/redis/go-redis/v9"
)

type CacheHelper struct {
	Cache *cache.Cache
}

func NewCacheHelper() *CacheHelper {
	// ring := redis.NewRing(&redis.RingOptions{
	// 	Addrs: map[string]string{
	// 		"server1": viper.GetString("redis.node1"),
	// 	},
	// })
	
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "localhost:6379",
		},
	})
	status, err := ring.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Redis connection was refused")
	}
	fmt.Println(status)
	mycache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(24, time.Hour),
	})
	return &CacheHelper{
		Cache: mycache,
	}
}

func (u *CacheHelper) Exist(ctx context.Context, key string) (res bool) {
	return u.Cache.Exists(ctx, key)
}

func (u *CacheHelper) Get(ctx context.Context, key string, value interface{}) (err error) {
	err = u.Cache.Get(ctx, key, value)
	return
}

func (u *CacheHelper) Set(ctx context.Context, key string, value interface{}) (err error) {
	err = u.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
	})
	return
}

func (u *CacheHelper) Delete(ctx context.Context, key string) (err error) {
	err = u.Cache.Delete(ctx, key)
	return
}


// func (u *CacheHelper)Exist(ctx context.Context,key string) (res bool){
// return
// }
