package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"xws/webook/internal/domain"

	"github.com/redis/go-redis/v9"
)

// var ErrKeyNotExist = errors.New("key 不存在")
var ErrKeyNotExist = redis.Nil

/*
type Cache interface{
	// 读用户信息
	GetUser(ctx context.Context, id int64) (domain.User, error)
	// 还有别的业务
}

// 最佳实践👇🏻
type CacheV1 interface{
	// 中间件团队去做
	Get(ctx context.Context, key string) (any, error)
}
// 底层基于redis/memecache来实现不同的缓存机制


type UserCache struct {
	cache CacheV1
}

func(u * UserCache) GetUser(ctx context.Context, id int64) (domain.User, error){

}
可是我们没有CacheV1
*/

type UserCache struct {
	// 传单机 Redis可以
	// 传cluster de Redis 也可以
	// client *redis.Client
	// ClusterClient *redis.ClusterClient
	client     redis.Cmdable
	expiration time.Duration
}

// 重要经验 💡
// A 用到了B, B一定是接口
// A 用到了B, B一定是A的字段
// A 用到了B, A 绝对不初始化B, 而是外面注入

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 只要 error 为nil, 就认为 缓存里有数据
// 如果没有数据,返回一个特定的error
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	// 数据不存在, err = redis.Nil
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err

}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *UserCache) key(id int64) string {
	// user:info:key
	// user_info_key
	return fmt.Sprintf("user:info:%d", id)
}
