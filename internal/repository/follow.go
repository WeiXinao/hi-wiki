package repository

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/xkit/slice"
	"time"
)

var ErrFollowedObjectNotExists = dao.ErrFollowedObjectNotExists

type FollowRepo interface {
	GetFollowInfo(ctx context.Context, uid []int64,
		currentUser int64, followedId []int64, typ uint8) (map[int64]domain.FollowInfo, error)
	Insert(ctx context.Context, typ uint8, followedId int64, userId int64) error
	GetFollowedArticleAndRepo(ctx context.Context, uid int64) ([]domain.FollowRecent, error)
}

type followRepo struct {
	dao dao.FollowDao
}

func (f *followRepo) GetFollowedArticleAndRepo(ctx context.Context,
	uid int64) ([]domain.FollowRecent, error) {
	followRecents, err := f.dao.GetFollowedArticleAndRepo(ctx, uid)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.FollowRecent, domain.FollowRecent](followRecents,
		func(idx int, src dao.FollowRecent) domain.FollowRecent {
			return f.FollowRecentEntityToDomain(src)
		}), nil
}

func (f *followRepo) Insert(ctx context.Context, typ uint8, followedId int64,
	userId int64) error {
	return f.dao.Insert(ctx, typ, followedId, userId)
}

func (f *followRepo) GetFollowInfo(ctx context.Context, uid []int64,
	currentUser int64, followedId []int64, typ uint8) (map[int64]domain.FollowInfo, error) {
	infoMap, err := f.dao.GetFollowInfo(ctx, uid, currentUser, followedId, typ)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.FollowInfo)
	for k, info := range infoMap {
		res[k] = f.FollowInfoEntityToDomain(info)
	}
	return res, nil
}

func (f *followRepo) FollowInfoEntityToDomain(info dao.FollowInfo) domain.FollowInfo {
	return domain.FollowInfo{
		FollowCount:   info.FollowCount,
		BeFollowCount: info.BeFollowCount,
		IsFollowed:    info.IsFollowed,
	}
}

func (f *followRepo) FollowRecentEntityToDomain(recent dao.FollowRecent) domain.FollowRecent {
	return domain.FollowRecent{
		UserName:   recent.UserName,
		Avatar:     "files/img/" + recent.Avatar,
		UniqueCode: recent.UniqueCode,
		Title:      recent.Title,
		CreateTime: time.UnixMilli(recent.CreateTime).Format("2006-01-02 03:04"),
		Typ:        domain.FollowType(recent.Typ).String(),
	}
}

func NewFollowRepository(dao dao.FollowDao) FollowRepo {
	return &followRepo{
		dao: dao,
	}
}
