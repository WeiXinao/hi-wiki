package service

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"github.com/gin-gonic/gin"
)

var ErrFollowedObjectNotExists = repository.ErrFollowedObjectNotExists

type FollowService interface {
	StatsFollowInfo(ctx context.Context, uids []int64,
		currentUser int64, followedIds []int64, typ uint8) (map[int64]domain.FollowInfo, error)
	Add(ctx context.Context, typ domain.FollowType, followedId int64, userId int64) error
	Recent(ctx *gin.Context, uid int64) ([]domain.FollowRecent, error)
}

type followService struct {
	repo repository.FollowRepo
}

func (f *followService) Recent(ctx *gin.Context, uid int64) ([]domain.FollowRecent, error) {
	return f.repo.GetFollowedArticleAndRepo(ctx, uid)
}

func (f *followService) Add(ctx context.Context, typ domain.FollowType, followedId int64,
	userId int64) error {
	return f.repo.Insert(ctx, typ.ToUint8(), followedId, userId)
}

func (f *followService) StatsFollowInfo(ctx context.Context, uids []int64, currentUser int64,
	followedIds []int64, typ uint8) (map[int64]domain.FollowInfo, error) {
	return f.repo.GetFollowInfo(ctx, uids, currentUser, followedIds, typ)
}

func NewFollowService(repo repository.FollowRepo) FollowService {
	return &followService{
		repo: repo,
	}
}
