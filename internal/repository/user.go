package repository

import (
	"context"
	"database/sql"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/cache"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"time"
)

var (
	ErrUserNotFound          = dao.ErrRecordNotFound
	ErrDuplicateUsername     = dao.ErrDuplicateUsername
	ErrPasswordHasBeenModify = dao.ErrPasswordHasBeenModify
	ErrFailModifyProfile     = dao.ErrFailModifyProfile
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	UpdateByIdAndPassword(ctx context.Context, id int64, oldPwd string, newPwd string) error
	UpdateById(ctx context.Context, avatarUrl string, desc string, uid int64) error
}

type cachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func (repo *cachedUserRepository) UpdateById(ctx context.Context, avatarUrl string, desc string, uid int64) error {
	err := repo.dao.UpdateById(ctx, avatarUrl, desc, uid)
	if err != nil {
		return err
	}
	time.AfterFunc(time.Second, func() {
		_ = repo.cache.Del(ctx, uid)
	})
	return repo.cache.Del(ctx, uid)
}

func (repo *cachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}

	user, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = repo.entityToDomain(user)

	go func() {
		_ = repo.cache.Set(ctx, u)
	}()

	return u, nil
}

func (repo *cachedUserRepository) UpdateByIdAndPassword(ctx context.Context, id int64, oldPwd string, newPwd string) error {
	return repo.dao.UpdateByIdAndPassword(ctx, id, oldPwd, newPwd)
}

func (repo *cachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(u))
}

func (repo *cachedUserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	u, err := repo.dao.FindByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(u), nil
}

func (repo *cachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:        u.Id,
		Username:  u.Username.String,
		Password:  u.Password,
		AvatarMd5: u.AvatarMd5,
		Profile:   u.Profile,
		CreateAt:  time.UnixMilli(u.CreatedAt),
		UpdatedAt: time.UnixMilli(u.UpdatedAt),
		DeletedAt: time.UnixMilli(u.DeletedAt),
	}
}

func (repo *cachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		BaseModel: dao.BaseModel{
			Id:        u.Id,
			CreatedAt: u.CreateAt.UnixMilli(),
			UpdatedAt: u.UpdatedAt.UnixMilli(),
			DeletedAt: u.DeletedAt.UnixMilli(),
		},
		Username: sql.NullString{
			String: u.Username,
			Valid:  u.Username != "",
		},
		Password:  u.Password,
		AvatarMd5: u.AvatarMd5,
		Profile:   u.Profile,
	}
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &cachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}
