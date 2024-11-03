package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	Del(ctx context.Context, id int64) error
}

type redisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *redisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := r.key(id)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (r *redisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := r.key(u.Id)
	return r.client.Set(ctx, key, val, r.expiration).Err()
}

func (r *redisUserCache) Del(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.key(id)).Err()
}

func (r *redisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func NewRedisUserCache(client redis.Cmdable) UserCache {
	return &redisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
