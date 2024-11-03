package repository

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/cache"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"time"
)

var ErrImageNotFound = dao.ErrImageNotFound

type FileRepository interface {
	InsertFile(ctx context.Context, file domain.File) (domain.File, error)
	InsertOwner(ctx context.Context, uid int64, md5 string) (domain.File, error)
	GetImageByMd5(ctx context.Context, md5Str string) (domain.File, error)
}

type fileRepository struct {
	dao   dao.FileDao
	cache cache.FileCache
}

func (f *fileRepository) GetImageByMd5(ctx context.Context, md5Str string) (domain.File, error) {
	fileDomain, err := f.cache.Get(ctx, md5Str)
	if err == nil {
		return fileDomain, nil
	}

	file, err := f.dao.GetImageByMd5(ctx, md5Str)
	if err != nil {
		return domain.File{}, err
	}
	fileDomain = f.FileEntityToDomain(file)

	go func() {
		_ = f.cache.Set(ctx, fileDomain)
	}()

	return fileDomain, nil
}

func (f *fileRepository) InsertOwner(ctx context.Context, uid int64, md5 string) (domain.File, error) {
	fileModel, err := f.dao.InsertOwner(ctx, uid, md5)
	if err != nil {
		return domain.File{}, err
	}
	return f.FileEntityToDomain(fileModel), err
}

func (f *fileRepository) InsertFile(ctx context.Context, file domain.File) (domain.File, error) {
	fileModel, err := f.dao.InsertFile(ctx, f.FileDomainToEntity(file))
	fileDomain := f.FileEntityToDomain(fileModel)
	// 预加载到缓存
	go func() {
		_ = f.cache.Set(ctx, fileDomain)
	}()
	return fileDomain, err
}

func (f *fileRepository) FileDomainToEntity(file domain.File) dao.File {
	return dao.File{
		BaseModel: dao.BaseModel{
			Id:        file.Id,
			CreatedAt: file.CreateAt.UnixMilli(),
			UpdatedAt: file.UpdatedAt.UnixMilli(),
			DeletedAt: file.DeletedAt.UnixMilli(),
		},
		Name:           file.Name,
		Typ:            file.Typ,
		Md5:            file.Md5,
		Url:            file.Url,
		Size:           file.Size,
		DirLevel:       file.DirLevel,
		RepoUniqueCode: file.RepoUniqueCode,
		UserId:         file.Owner.Id,
	}
}

func (f *fileRepository) FileEntityToDomain(file dao.File) domain.File {
	return domain.File{
		Id:             file.Id,
		Name:           file.Name,
		Typ:            file.Typ,
		Md5:            file.Md5,
		Url:            file.Url,
		Size:           file.Size,
		DirLevel:       file.DirLevel,
		RepoUniqueCode: file.RepoUniqueCode,
		Owner: domain.Owner{
			Id: file.Id,
		},
		CreateAt:  time.UnixMilli(file.CreatedAt),
		UpdatedAt: time.UnixMilli(file.UpdatedAt),
		DeletedAt: time.UnixMilli(file.DeletedAt),
	}
}

func NewFileRepository(dao dao.FileDao, cache cache.FileCache) FileRepository {
	return &fileRepository{
		dao:   dao,
		cache: cache,
	}
}
