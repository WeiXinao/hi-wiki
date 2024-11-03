package repository

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/xkit/slice"
	"time"
)

var (
	ErrGoodHistoryNotFound = dao.ErrRecordNotFound
	ErrHasGood             = dao.ErrHasGood
)

type GoodRepository interface {
	GetByUserIdAndArticleUniqueCode(ctx context.Context, uid int64, uniqueCodes []string) ([]domain.GoodHistory, error)
	InsertGood(ctx context.Context, code string, uid int64) error
}

type goodRepository struct {
	dao dao.GoodDao
}

func (g *goodRepository) InsertGood(ctx context.Context, code string, uid int64) error {
	return g.dao.InsertGood(ctx, code, uid)
}

func (g *goodRepository) GetByUserIdAndArticleUniqueCode(ctx context.Context, uid int64,
	uniqueCodes []string) ([]domain.GoodHistory, error) {
	goodHistories, err := g.dao.GetByUserIdAndArticleUniqueCode(ctx, uid, uniqueCodes)
	return slice.Map[dao.GoodHistory, domain.GoodHistory](goodHistories,
		func(idx int, src dao.GoodHistory) domain.GoodHistory {
			return g.entityToDomain(src)
		}), err
}

func (g *goodRepository) entityToDomain(goodHistory dao.GoodHistory) domain.GoodHistory {
	return domain.GoodHistory{
		Id:                goodHistory.Id,
		UserId:            goodHistory.UserId,
		ArticleUniqueCode: goodHistory.ArticleUniqueCode,
		CreateAt:          time.UnixMilli(goodHistory.CreatedAt),
		UpdatedAt:         time.UnixMilli(goodHistory.UpdatedAt),
		DeletedAt:         time.UnixMilli(goodHistory.DeletedAt),
	}
}

func NewGoodRepository(dao dao.GoodDao) GoodRepository {
	return &goodRepository{
		dao: dao,
	}
}
