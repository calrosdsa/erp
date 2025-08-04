package cache

import (
	"erp/internal/app/connection"
	"erp/internal/app/domain"
	"fmt"

	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	cache *ristretto.Cache
	conn *connection.Connection
}

func NewCache(
	conn *connection.Connection,
) *Cache {
	// ring := redis.NewRing(&redis.RingOptions{
	// 	Addrs: map[string]string{
	// 		"server1": viper.GetString("redis.node1"),
	// 	},
	// })
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return &Cache{
		cache: cache,
		conn: conn,
	}
}

// func (u *Cache) Get

// func (u *Cache) Exist(ctx context.Context, key string) (res bool) {
// 	return u.Cache.
// }

func (u *Cache) Del(opts ...Option) {
	o := Options.Apply(opts...)
	key := o.KeyEntity
	if o.ID != "" {
		key = fmt.Sprintf("%s_%s",key,o.ID)
	}
	u.del(key)
}

func (u *Cache) Get(opts ...Option) ( bool) {
	o := Options.Apply(opts...)
	key := o.KeyEntity
	if o.ID != "" {
		key = fmt.Sprintf("%s_%s",key,o.ID)
	}
	if res,found := u.get(key);found {
		fmt.Println("VALUE",res)
		o.TypeAssertion(res)
		fmt.Println("VALUE 2",o.Data)
		return found
	}
	return false
}

func (u *Cache) Set(opts ...Option) bool {
	o := Options.Apply(opts...)
	key := o.KeyEntity
	if o.ID != "" {
		key = fmt.Sprintf("%s_%s",key,o.ID)
	}
	ok := u.set(key,o.Data)
	return ok
}

func (u *Cache) GetEntity(opts ...Option) error{
	o := Options.Apply(opts...)
	key := o.KeyEntity
	if o.ID != "" {
		key = fmt.Sprintf("%s_%s",key,o.ID)
	}
	if res,found := u.get(key);found {
		o.TypeAssertion(res)
		return nil
	}else {
		err := u.conn.Db.Where(o.Data).First(o.Data).Error
		if err != nil {
			return domain.UNEXPECTED_ERROR
		}
		ok := u.set(key,o.Data)
		if ok {
			fmt.Println("getting from db")
		}
		return nil
	}
}


func (u *Cache) get(key string) (interface{},bool) {
	return u.cache.Get(key)
}

func (u *Cache) set(key string, value interface{}) (bool) {
	res := u.cache.Set(key,value,0)
	u.cache.Wait()
	return res
}

func (u *Cache) del(key string) {
	u.cache.Del(key)
}


// func (u *Cache)Exist(ctx context.Context,key string) (res bool){
// return
// }
