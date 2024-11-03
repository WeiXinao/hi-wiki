package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, article Article) error
	UpdateByUniqueCode(ctx context.Context, article Article) error
	GetPageByUserId(ctx context.Context, uid int64, offset int,
		size int) ([]Article, int64, error)
	GetByUniqueCodeAndUserId(ctx context.Context, uniqueCode string, uid int64) (Article, error)
	GetPages(ctx context.Context, offset int, size int) ([]Article, int64, error)
	GetPagesByHot(ctx context.Context, offset int, size int) ([]Article, int64, error)
}

type articleMinioDao struct {
	db *gorm.DB
}

func (g *articleMinioDao) GetPagesByHot(ctx context.Context, offset int,
	size int) ([]Article, int64, error) {
	var count int64
	var articles []Article
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Article{}).
			Where("state = ?", 0).
			Order("like_cnt DESC").
			Count(&count).
			Error
		if err != nil {
			return err
		}

		return tx.Where("state = ?", 0).
			Order("like_cnt DESC").
			Offset((offset - 1) * size).
			Limit(size).
			Preload("User").
			Preload("RepoCate").
			Find(&articles).Error
	})
	if err != nil {
		return nil, 0, err
	}
	return articles, count, nil
}

func (g *articleMinioDao) GetPages(ctx context.Context, offset int,
	size int) ([]Article, int64, error) {
	var count int64
	var articles []Article
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Article{}).Count(&count).Error
		if err != nil {
			return err
		}

		return tx.
			Offset((offset - 1) * size).
			Limit(size).
			Preload("User").
			Preload("RepoCate").
			Find(&articles).Error
	})
	if err != nil {
		return nil, 0, err
	}
	return articles, count, nil
}

func (g *articleMinioDao) GetByUniqueCodeAndUserId(ctx context.Context, uniqueCode string, uid int64) (Article, error) {
	var art Article
	err := g.db.WithContext(ctx).
		Where("unique_code = ?", uniqueCode).
		Preload("Repo").
		Preload("User").
		First(&art).
		Error
	return art, err
}

func (g *articleMinioDao) GetPageByUserId(ctx context.Context, uid int64, offset int,
	size int) ([]Article, int64, error) {
	var count int64
	var articles []Article
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Article{}).Where("user_id = ?", uid).Count(&count).Error
		if err != nil {
			return err
		}

		return tx.Where("user_id = ?", uid).
			Offset((offset - 1) * size).
			Limit(size).
			Preload("User").
			Preload("RepoCate").
			Find(&articles).Error
	})
	if err != nil {
		return nil, 0, err
	}
	return articles, count, nil
}

func (g *articleMinioDao) UpdateByUniqueCode(ctx context.Context, article Article) error {
	return g.db.WithContext(ctx).
		Model(&article).
		Where("unique_code = ?", article.UniqueCode).
		Updates(map[string]any{
			"title": article.Title,
		}).Error
}

func (g *articleMinioDao) Insert(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.CreatedAt = now
	article.UpdatedAt = now
	article.DeletedAt = 0
	return g.db.WithContext(ctx).Create(&article).Error
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &articleMinioDao{
		db: db,
	}
}

type Article struct {
	BaseModel
	Title string
	//Content        string
	//PureContent    string
	LikeCnt        int64
	UniqueCode     string `gorm:"unique"`
	RepoUniqueCode string `gorm:"type:varchar(191)"`
	UserId         int64
	Cate           int64
	State          uint8
	Private        bool
	RepoCate       RepoCate `gorm:"foreignKey:Cate"`
	Repo           Repo     `gorm:"foreignKey:RepoUniqueCode;references:UniqueCode"`
	User           User     `gorm:"foreignKey:UserId"`
}
