package service

import (
	"context"
	"errors"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/xkit/slice"
)

var ErrHasGood = dao.ErrHasGood

type GoodService interface {
	IsLiked(ctx context.Context, uid int64, uniqueCodes []string) (map[string]bool, error)
	GoodArticle(ctx context.Context, code string, uid int64) error
}

type goodService struct {
	repo repository.GoodRepository
}

func (g *goodService) GoodArticle(ctx context.Context, code string, uid int64) error {
	return g.repo.InsertGood(ctx, code, uid)
}

func (g *goodService) IsLiked(ctx context.Context, uid int64, uniqueCodes []string) (map[string]bool, error) {
	goodHistories, err := g.repo.GetByUserIdAndArticleUniqueCode(ctx, uid, uniqueCodes)
	if errors.Is(err, repository.ErrGoodHistoryNotFound) {
		return nil, repository.ErrGoodHistoryNotFound
	}
	if err != nil {
		return nil, err
	}
	res := make(map[string]bool)
	for _, uniqueCode := range uniqueCodes {
		res[uniqueCode] = slice.ContainsFunc[domain.GoodHistory](goodHistories, func(src domain.GoodHistory) bool {
			return src.ArticleUniqueCode == uniqueCode && src.UserId == uid
		})
	}
	return res, nil
}

func NewGoodService(repo repository.GoodRepository) GoodService {
	return &goodService{
		repo: repo,
	}
}
