package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrArticleNotFound = ErrRecordNotFound
	ErrHasGood         = errors.New("已经点过赞了")
)

type GoodDao interface {
	GetByUserIdAndArticleUniqueCode(ctx context.Context, uid int64, codes []string) ([]GoodHistory, error)
	InsertGood(ctx context.Context, code string, uid int64) error
}

type GormGoodDao struct {
	db *gorm.DB
}

func (g *GormGoodDao) InsertGood(ctx context.Context, code string, uid int64) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&Article{}).
			Where("unique_code = ?", code).
			Update("like_cnt", gorm.Expr("like_cnt + ?", 1))
		err := result.Error
		if err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			return ErrArticleNotFound
		}
		now := time.Now().UnixMilli()
		err = tx.Create(&GoodHistory{
			BaseModel: BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: 0,
			},
			ArticleUniqueCode: code,
			UserId:            uid,
		}).Error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			const uniqueConflictsErrNo uint16 = 1062
			if mysqlErr.Number == uniqueConflictsErrNo {
				return ErrHasGood
			}
		}
		return nil
	})
}

func (g *GormGoodDao) GetByUserIdAndArticleUniqueCode(ctx context.Context, uid int64,
	codes []string) ([]GoodHistory, error) {
	var goodHistories []GoodHistory
	err := g.db.WithContext(ctx).
		Model(&GoodHistory{}).
		Where("user_id = ? AND article_unique_code IN ?", uid, codes).
		Find(&goodHistories).Error
	if err != nil {
		return nil, err
	}
	return goodHistories, err
}

func NewGoodDao(db *gorm.DB) GoodDao {
	return &GormGoodDao{
		db: db,
	}
}

type GoodHistory struct {
	BaseModel
	UserId            int64  `gorm:"unique_index:uk_uid_article_unique_code"`
	ArticleUniqueCode string `gorm:"unique_index:uk_uid_article_unique_code;type:varchar(191)"`
}
