package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"strings"
)

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) error
	ListSelf(ctx context.Context, uid int64, offset int,
		size int) ([]domain.Article, int64, error)
	Detail(ctx context.Context, uniqueCode string, uid int64) (domain.Article, error)
	ListAll(ctx context.Context, offset int,
		size int, isHot bool) ([]domain.Article, int64, error)
}

type articleService struct {
	repo     repository.ArticleRepository
	repoRepo repository.RepoRepository
}

func (a *articleService) ListAll(ctx context.Context, offset int,
	size int, isHot bool) ([]domain.Article, int64, error) {
	if isHot {
		return a.repo.GetPagesByHot(ctx, offset, size)
	}
	return a.repo.GetPages(ctx, offset, size)
}

func (a *articleService) Detail(ctx context.Context, uniqueCode string, uid int64) (domain.Article, error) {
	return a.repo.GetByUniqueCodeAndUserId(ctx, uniqueCode, uid)
}

func (a *articleService) ListSelf(ctx context.Context, uid int64, offset int,
	size int) ([]domain.Article, int64, error) {
	return a.repo.GetPageByUserId(ctx, uid, offset, size)
}

func (a *articleService) Edit(ctx context.Context, article domain.Article) error {
	article.Cate.Id = 1
	repo, err := a.repoRepo.GetByUniqueCode(ctx, article.RepoUniqueCode)
	if err != nil {
		return err
	}
	if repo.Status.IsPrivate() {
		article.Private = true
	}
	if len(strings.TrimSpace(article.UniqueCode)) == 0 {
		uniqueCode, err := a.generateRandomString(32)
		if err != nil {
			return err
		}
		article.UniqueCode = uniqueCode
		err = a.repo.Insert(ctx, article)
		if err != nil {
			return err
		}
	} else {
		err := a.repo.Update(ctx, article)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *articleService) generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewArticleService(repo repository.ArticleRepository,
	repoRepo repository.RepoRepository) ArticleService {
	return &articleService{
		repo:     repo,
		repoRepo: repoRepo,
	}
}
