package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/hi-wiki/internal/repository/oss"
	"github.com/WeiXinao/hi-wiki/pkg/logger"
	"github.com/WeiXinao/xkit/slice"
	"github.com/minio/minio-go/v7"
	"io"
	"time"
)

type ArticleRepository interface {
	Insert(ctx context.Context, article domain.Article) error
	Update(ctx context.Context, article domain.Article) error
	GetPageByUserId(ctx context.Context, uid int64, offset int,
		size int) ([]domain.Article, int64, error)
	GetPages(ctx context.Context, offset int, size int) ([]domain.Article, int64, error)
	GetByUniqueCodeAndUserId(ctx context.Context, uniqueCode string, uid int64) (domain.Article, error)
	GetPagesByHot(ctx context.Context, offset int, size int) ([]domain.Article, int64, error)
}

type articleRepository struct {
	dao                dao.ArticleDao
	oss                oss.OSS
	contentKeyTmpl     string
	pureContentKeyTmpl string
	textPutOpt         minio.PutObjectOptions
	htmlPutOpt         minio.PutObjectOptions
	l                  logger.Logger
}

func (a *articleRepository) GetPagesByHot(ctx context.Context, offset int, size int) ([]domain.Article, int64, error) {
	arts, count, err := a.dao.GetPagesByHot(ctx, offset, size)
	if err != nil {
		return nil, 0, err
	}
	artsDomains := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		art := a.entityToDomain(src)

		object, err := a.oss.GetObject(ctx, fmt.Sprintf(a.contentKeyTmpl, src.UniqueCode), minio.GetObjectOptions{})
		if err != nil {
			a.l.Error("获取 oss object 失败", logger.Error(err))
		}
		defer object.Close()

		bytes, err := io.ReadAll(object)
		if err != nil {
			a.l.Error("读取 object content 失败", logger.Error(err))
		}
		art.Content = string(bytes)

		object, err = a.oss.GetObject(ctx, fmt.Sprintf(a.pureContentKeyTmpl, src.UniqueCode), minio.GetObjectOptions{})
		if err != nil {
			a.l.Error("获取 oss object 失败", logger.Error(err))
		}
		bytes, err = io.ReadAll(object)
		if err != nil {
			a.l.Error("读取 object pureContent 失败", logger.Error(err))
		}
		art.PureContent = string(bytes)

		art.GenDesc()
		return art
	})
	return artsDomains, count, nil
}

func (a *articleRepository) GetPages(ctx context.Context, offset int, size int) ([]domain.Article, int64, error) {
	arts, count, err := a.dao.GetPages(ctx, offset, size)
	if err != nil {
		return nil, 0, err
	}
	artsDomains := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		art := a.entityToDomain(src)

		object, err := a.oss.GetObject(ctx, fmt.Sprintf(a.contentKeyTmpl, src.UniqueCode), minio.GetObjectOptions{})
		if err != nil {
			a.l.Error("获取 oss object 失败", logger.Error(err))
		}
		defer object.Close()

		bytes, err := io.ReadAll(object)
		if err != nil {
			a.l.Error("读取 object content 失败", logger.Error(err))
		}
		art.Content = string(bytes)

		object, err = a.oss.GetObject(ctx, fmt.Sprintf(a.pureContentKeyTmpl, src.UniqueCode), minio.GetObjectOptions{})
		if err != nil {
			a.l.Error("获取 oss object 失败", logger.Error(err))
		}
		bytes, err = io.ReadAll(object)
		if err != nil {
			a.l.Error("读取 object pureContent 失败", logger.Error(err))
		}
		art.PureContent = string(bytes)

		art.GenDesc()
		return art
	})
	return artsDomains, count, nil
}

func (a *articleRepository) GetByUniqueCodeAndUserId(ctx context.Context, uniqueCode string, uid int64) (domain.Article, error) {
	art, err := a.dao.GetByUniqueCodeAndUserId(ctx, uniqueCode, uid)
	artDomain := a.entityToDomain(art)
	object, err := a.oss.GetObject(ctx, fmt.Sprintf(a.contentKeyTmpl, artDomain.UniqueCode), minio.GetObjectOptions{})
	if err != nil {
		a.l.Error("获取 oss object 失败", logger.Error(err))
	}
	defer object.Close()
	bytes, err := io.ReadAll(object)
	if err != nil {
		return domain.Article{}, nil
	}
	artDomain.Content = string(bytes)
	return artDomain, err
}

func (a *articleRepository) GetPageByUserId(ctx context.Context, uid int64, offset int,
	size int) ([]domain.Article, int64, error) {
	arts, count, err := a.dao.GetPageByUserId(ctx, uid, offset, size)
	if err != nil {
		return nil, 0, err
	}
	artsDomains := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		art := a.entityToDomain(src)

		object, err := a.oss.GetObject(ctx, fmt.Sprintf(a.contentKeyTmpl, src.UniqueCode), minio.GetObjectOptions{})
		if err != nil {
			a.l.Error("获取 oss object 失败", logger.Error(err))
		}
		defer object.Close()

		bytes, err := io.ReadAll(object)
		if err != nil {
			a.l.Error("读取 object content 失败", logger.Error(err))
		}
		art.Content = string(bytes)

		object, err = a.oss.GetObject(ctx, fmt.Sprintf(a.pureContentKeyTmpl, src.UniqueCode), minio.GetObjectOptions{})
		if err != nil {
			a.l.Error("获取 oss object 失败", logger.Error(err))
		}
		bytes, err = io.ReadAll(object)
		if err != nil {
			a.l.Error("读取 object pureContent 失败", logger.Error(err))
		}
		art.PureContent = string(bytes)

		art.GenDesc()
		return art
	})
	return artsDomains, count, nil
}

func (a *articleRepository) Update(ctx context.Context, article domain.Article) error {
	var (
		preContent     = article.Content
		prePure        = article.PureContent
		contentKey     = fmt.Sprintf(a.contentKeyTmpl, article.UniqueCode)
		pureContentKey = fmt.Sprintf(a.pureContentKeyTmpl, article.UniqueCode)
		content        = article.Content
		pureContent    = article.PureContent
	)
	_, err := a.oss.PutObject(ctx, contentKey,
		bytes.NewBufferString(content), int64(len(content)), a.htmlPutOpt)
	if err != nil {
		return err
	}
	_, err = a.oss.PutObject(ctx, pureContentKey,
		bytes.NewBufferString(pureContent), int64(len(pureContent)), a.textPutOpt)
	if err != nil {
		go func() {
			_, err = a.oss.PutObject(ctx, contentKey,
				bytes.NewBufferString(preContent), int64(len(preContent)), a.textPutOpt)
		}()
		return err
	}
	err = a.dao.UpdateByUniqueCode(ctx, a.domainToEntity(article))
	if err != nil {
		go func() {
			_, err = a.oss.PutObject(ctx, contentKey,
				bytes.NewBufferString(preContent), int64(len(preContent)), a.htmlPutOpt)
			_, err = a.oss.PutObject(ctx, pureContentKey,
				bytes.NewBufferString(prePure), int64(len(prePure)), a.textPutOpt)
		}()
		return err
	}
	return nil
}

func (a *articleRepository) Insert(ctx context.Context, article domain.Article) error {
	content := article.Content
	pureContent := article.PureContent
	contentKey := fmt.Sprintf(a.contentKeyTmpl, article.UniqueCode)
	pureContentKey := fmt.Sprintf(a.pureContentKeyTmpl, article.UniqueCode)
	_, err := a.oss.PutObject(ctx, contentKey, bytes.NewBufferString(content), int64(len(content)), a.htmlPutOpt)
	if err != nil {
		return err
	}
	_, err = a.oss.PutObject(ctx, pureContentKey,
		bytes.NewBufferString(pureContent), int64(len(pureContent)), a.textPutOpt)
	if err != nil {
		go func() {
			newCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = a.oss.RemoveObject(newCtx, contentKey, minio.RemoveObjectOptions{})
		}()
		return err
	}
	err = a.dao.Insert(ctx, a.domainToEntity(article))
	if err != nil {
		newCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = a.oss.RemoveObject(newCtx, contentKey, minio.RemoveObjectOptions{})
		_ = a.oss.RemoveObject(newCtx, pureContentKey, minio.RemoveObjectOptions{})
		return err
	}
	return nil
}

func (a *articleRepository) domainToEntity(article domain.Article) dao.Article {
	return dao.Article{
		BaseModel: dao.BaseModel{
			Id:        article.Id,
			CreatedAt: article.CreateAt.UnixMilli(),
			UpdatedAt: article.UpdatedAt.UnixMilli(),
			DeletedAt: article.DeletedAt.UnixMilli(),
		},
		Title:          article.Title,
		LikeCnt:        article.LikeCnt,
		UniqueCode:     article.UniqueCode,
		RepoUniqueCode: article.RepoUniqueCode,
		UserId:         article.Author.Id,
		Cate:           article.Cate.Id,
		State:          article.State,
		Private:        article.Private,
	}
}

func (a *articleRepository) entityToDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:             art.Id,
		Title:          art.Title,
		LikeCnt:        art.LikeCnt,
		UniqueCode:     art.UniqueCode,
		RepoUniqueCode: art.RepoUniqueCode,
		State:          art.State,
		Private:        art.Private,
		CreateAt:       time.UnixMilli(art.CreatedAt),
		UpdatedAt:      time.UnixMilli(art.UpdatedAt),
		DeletedAt:      time.UnixMilli(art.DeletedAt),
		Author: domain.Author{
			Id:        art.User.Id,
			Name:      art.User.Username.String,
			AvatarMd5: art.User.AvatarMd5,
			Profile:   art.User.Profile,
		},
		Repo: domain.Repo{
			Id:         art.Repo.Id,
			Name:       art.Repo.Name,
			Desc:       art.Repo.Desc,
			Status:     domain.RepoStatus(art.Repo.Status),
			UniqueCode: art.Repo.UniqueCode,
			State:      art.Repo.State,
		},
		Cate: domain.Cate{
			Id:   art.RepoCate.Id,
			Name: art.RepoCate.Name,
		},
	}
}

func NewArticleRepository(dao dao.ArticleDao, oss oss.OSS, l logger.Logger) ArticleRepository {
	return &articleRepository{
		dao:                dao,
		oss:                oss,
		contentKeyTmpl:     "article/%s/content.html",
		pureContentKeyTmpl: "article/%s/purecontent.txt",
		textPutOpt: minio.PutObjectOptions{
			ContentType: "text/plain;charset=utf-8",
		},
		htmlPutOpt: minio.PutObjectOptions{
			ContentType: "text/html;charset=utf-8",
		},
		l: l,
	}
}
