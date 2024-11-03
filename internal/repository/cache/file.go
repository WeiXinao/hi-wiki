package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type FileCache interface {
	Get(ctx context.Context, md5 string) (domain.File, error)
	Set(ctx context.Context, file domain.File) error
}

type redisFileCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *redisFileCache) Get(ctx context.Context, md5 string) (domain.File, error) {
	key := r.key(md5)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.File{}, err
	}
	var f domain.File
	err = json.Unmarshal(val, &f)
	return f, err
}

func (r *redisFileCache) Set(ctx context.Context, file domain.File) error {
	val, err := json.Marshal(file)
	if err != nil {
		return err
	}
	key := r.key(file.Md5)
	return r.client.Set(ctx, key, val, r.expiration).Err()
}

func (r *redisFileCache) key(md5 string) string {
	return fmt.Sprintf("file:md5:%s", md5)
}

func NewRedisFileCache(client redis.Cmdable) FileCache {
	return &redisFileCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
