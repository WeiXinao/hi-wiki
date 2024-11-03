package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
)

type RepoService interface {
	Create(ctx context.Context, repo domain.Repo) error
	List(ctx context.Context, uid int64, isDoc bool) ([]domain.Repo, error)
	ListByTeamId(ctx context.Context, groupId int64) ([]domain.Repo, error)
	ListHot(ctx context.Context, isDoc bool, isHot bool) ([]domain.Repo, error)
}

type repoService struct {
	repo repository.RepoRepository
}

func (r *repoService) ListHot(ctx context.Context, isDoc bool, isHot bool) ([]domain.Repo, error) {
	return r.repo.GetHotDocRepo(ctx, isDoc, isHot)
}

func (r *repoService) ListByTeamId(ctx context.Context, groupId int64) ([]domain.Repo, error) {
	return r.repo.GetByGroupId(ctx, groupId)
}

func (r *repoService) List(ctx context.Context, uid int64, isDoc bool) ([]domain.Repo, error) {
	return r.repo.GetByUserId(ctx, uid, isDoc)
}

func (r *repoService) Create(ctx context.Context, repo domain.Repo) error {
	uniqueCode, err := r.generateRandomString(32)
	if err != nil {
		return err
	}
	repo.UniqueCode = uniqueCode
	return r.repo.InsertRepo(ctx, repo)
}

func (r *repoService) generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func NewRepoService(repo repository.RepoRepository) RepoService {
	return &repoService{
		repo: repo,
	}
}
