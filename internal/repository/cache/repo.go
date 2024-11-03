package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RepoCache interface {
	Set(ctx context.Context, repo domain.Repo) error
}

type redisRepoCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *redisRepoCache) Set(ctx context.Context, repo domain.Repo) error {
	val, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	key := r.key(repo.UniqueCode)
	return r.client.Set(ctx, key, val, r.expiration).Err()
}

func (r *redisRepoCache) key(uniqueCode string) string {
	return fmt.Sprintf("repo:info:%s", uniqueCode)
}

func NewRedisRepoCache(client redis.Cmdable) RepoCache {
	return &redisRepoCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
