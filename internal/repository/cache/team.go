package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type TeamCache interface {
	GetTeam(ctx context.Context, code string) (domain.Team, error)
	SetTeam(ctx context.Context, team domain.Team) error
	AddTeamMembers(ctx context.Context, code string, members ...domain.TeamMember) error
	GetTeamMembers(ctx context.Context, code string) ([]domain.TeamMember, error)
	DelTeamMember(ctx context.Context, code string, uid int64) error
}

type RedisTeamCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *RedisTeamCache) GetTeam(ctx context.Context, code string) (domain.Team, error) {
	key := r.infoKey(code)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.Team{}, err
	}
	var t domain.Team
	err = json.Unmarshal(val, &t)
	return t, err
}

func (r *RedisTeamCache) SetTeam(ctx context.Context, team domain.Team) error {
	val, err := json.Marshal(team)
	if err != nil {
		return err
	}
	key := r.infoKey(team.UniqueCode)
	return r.client.Set(ctx, key, val, r.expiration).Err()
}

func (r *RedisTeamCache) infoKey(uniqueCode string) string {
	return fmt.Sprintf("team:info:%s", uniqueCode)
}

func (r *RedisTeamCache) AddTeamMembers(ctx context.Context, code string, members ...domain.TeamMember) error {
	key := r.memberKey(code)
	vals := make([]any, 0, 2*len(members))
	for _, member := range members {
		bytes, err := json.Marshal(member)
		if err != nil {
			continue
		}
		vals = append(vals, strconv.FormatInt(member.UserId, 10))
		vals = append(vals, bytes)
	}
	return r.client.HMSet(ctx, key, vals...).Err()
}

func (r *RedisTeamCache) GetTeamMembers(ctx context.Context, code string) ([]domain.TeamMember, error) {
	key := r.memberKey(code)
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if result == 0 {
		return nil, redis.Nil
	}
	members, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res := make([]domain.TeamMember, 0, len(members))
	tm := domain.TeamMember{}
	for _, number := range members {
		err := json.Unmarshal([]byte(number), &tm)
		if err != nil {
			continue
		}
		res = append(res, tm)
	}
	return res, nil
}

func (r *RedisTeamCache) DelTeamMember(ctx context.Context, code string, uid int64) error {
	key := r.memberKey(code)
	return r.client.HDel(ctx, key, strconv.FormatInt(uid, 10)).Err()
}

func (r *RedisTeamCache) memberKey(uniqueCode string) string {
	return fmt.Sprintf("team:members:%s", uniqueCode)
}

func NewRedisTeamCache(client redis.Cmdable) TeamCache {
	return &RedisTeamCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}
