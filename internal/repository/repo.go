package repository

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/cache"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/xkit/slice"
	"time"
)

type RepoRepository interface {
	InsertRepo(ctx context.Context, repo domain.Repo) error
	GetByUserId(ctx context.Context, uid int64, isDoc bool) ([]domain.Repo, error)
	GetByGroupId(ctx context.Context, groupId int64) ([]domain.Repo, error)
	GetByUniqueCode(ctx context.Context, uniqueCode string) (domain.Repo, error)
	GetHotDocRepo(ctx context.Context, isDoc bool, isHot bool) ([]domain.Repo, error)
}

type repoRepository struct {
	dao   dao.RepoDao
	cache cache.RepoCache
}

func (r *repoRepository) GetHotDocRepo(ctx context.Context, isDoc bool, isHot bool) ([]domain.Repo, error) {
	repos, err := r.dao.GetHotDocRepo(ctx, isDoc, isHot)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Repo, domain.Repo](repos, func(idx int, src dao.Repo) domain.Repo {
		return r.entityToDomain(src)
	}), nil
}

func (r *repoRepository) GetByUniqueCode(ctx context.Context, uniqueCode string) (domain.Repo, error) {
	repo, err := r.dao.GetByUniqueCode(ctx, uniqueCode)
	return r.entityToDomain(repo), err
}

func (r *repoRepository) GetByGroupId(ctx context.Context, groupId int64) ([]domain.Repo, error) {
	repos, err := r.dao.GetByGroupId(ctx, groupId)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Repo, domain.Repo](repos, func(idx int, src dao.Repo) domain.Repo {
		return r.entityToDomain(src)
	}), nil
}

func (r *repoRepository) GetByUserId(ctx context.Context, uid int64, isDoc bool) ([]domain.Repo, error) {
	repos, err := r.dao.GetByUserId(ctx, uid, isDoc)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Repo, domain.Repo](repos, func(idx int, src dao.Repo) domain.Repo {
		return r.entityToDomain(src)
	}), nil
}

func (r *repoRepository) InsertRepo(ctx context.Context, repo domain.Repo) error {
	err := r.dao.InsertRepo(ctx, r.domainToEntity(repo))
	if err != nil {
		return err
	}
	// 插入成功后预加载到缓存
	go func() {
		_ = r.cache.Set(ctx, repo)
	}()
	return nil
}

func (r *repoRepository) entityToDomain(repo dao.Repo) domain.Repo {
	return domain.Repo{
		Id:   repo.Id,
		Name: repo.Name,
		Desc: repo.Desc,
		Cate: domain.Cate{
			Id:   repo.RepoCate.Id,
			Name: repo.RepoCate.Name,
		},
		Status:     domain.RepoStatus(repo.Status),
		UniqueCode: repo.UniqueCode,
		State:      repo.State,
		UserGroup: domain.UserGroup{
			Id: repo.TeamId,
		},
		Creator: domain.Creator{
			Id: repo.UserId,
		},
		CreateAt:  time.UnixMilli(repo.CreatedAt),
		UpdatedAt: time.UnixMilli(repo.UpdatedAt),
		DeletedAt: time.UnixMilli(repo.DeletedAt),
	}
}

func (r *repoRepository) domainToEntity(repo domain.Repo) dao.Repo {
	return dao.Repo{
		BaseModel: dao.BaseModel{
			Id:        repo.Id,
			CreatedAt: repo.CreateAt.UnixMilli(),
			UpdatedAt: repo.UpdatedAt.UnixMilli(),
			DeletedAt: repo.DeletedAt.UnixMilli(),
		},
		Name:       repo.Name,
		Desc:       repo.Desc,
		Cate:       repo.Cate.Id,
		Status:     repo.Status.ToUint8(),
		UniqueCode: repo.UniqueCode,
		State:      repo.State,
		TeamId:     repo.UserGroup.Id,
		UserId:     repo.Creator.Id,
	}
}

func NewRepoRepository(dao dao.RepoDao, cache cache.RepoCache) RepoRepository {
	return &repoRepository{
		dao:   dao,
		cache: cache,
	}
}
