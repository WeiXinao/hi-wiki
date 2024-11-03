package dao

import (
	"context"
	"errors"
	"github.com/WeiXinao/hi-wiki/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrDeletedBookCateNotFound = errors.New("被删除的分类未找到")

type BookDao interface {
	InsertOrUpdate(ctx context.Context, id int64, name string) error
	List(ctx context.Context) ([]BookCate, error)
	Del(ctx context.Context, id int64) error
	Insert(ctx context.Context, bookName string, bookUrl string, bookAvatarUrl string,
		bookCateId int64, bookUserId int64, bookMd5 string, avatarMd5 string) error
	Get(ctx context.Context, cate int64, kw string, rank string, page int, size int) ([]Book, error)
	UpdateAndGet(ctx context.Context, id int64) (Book, error)
}

type gormBookDao struct {
	db *gorm.DB
	l  logger.Logger
}

func (g *gormBookDao) UpdateAndGet(ctx context.Context, id int64) (Book, error) {
	var (
		book Book
		now  = time.Now().UnixMilli()
	)

	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", id).Find(&book).Error
		if err != nil {
			return err
		}
		return tx.Model(&Book{}).Where("id = ? and download = ?", id, book.Download).
			Updates(map[string]any{
				"download":   gorm.Expr("download + ?", 1),
				"updated_at": now,
			}).Error
	})
	if err != nil {
		return Book{}, err
	}
	book.Download += 1
	return book, err
}

func (g *gormBookDao) Get(ctx context.Context, cate int64, kw string, rank string,
	page int, size int) ([]Book, error) {
	var books []Book
	db := g.db.WithContext(ctx).Model(&Book{})
	switch cate {
	case 0:
	default:
		db.Where("book_cate_id", cate)
	}
	if rank == "hot" {
		db.Order("download DESC")
	}
	err := db.Limit(size).Offset((page - 1) * size).Preload("BookCate").Find(&books).Error
	g.l.Debug("get sql:", logger.Int("size", size), logger.Int("page", page))

	if err != nil {
		return nil, err
	}
	return books, nil
}

func (g *gormBookDao) Insert(ctx context.Context, bookName string, bookUrl string, bookAvatarUrl string,
	bookCateId int64, bookUserId int64, bookMd5 string, avatarMd5 string) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).
		Create(&Book{
			BaseModel: BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: 0,
			},
			BookName:      bookName,
			BookUrl:       bookUrl,
			BookAvatarUrl: bookAvatarUrl,
			BookCateId:    bookCateId,
			BookUserId:    bookUserId,
			BookMd5:       bookMd5,
			AvatarMd5:     avatarMd5,
		}).Error
}

func (g *gormBookDao) Del(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	res := g.db.WithContext(ctx).
		Model(&BookCate{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"deleted_at": now,
		})
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return ErrDeletedBookCateNotFound
	}
	return err
}

func (g *gormBookDao) List(ctx context.Context) ([]BookCate, error) {
	var bookCates []BookCate
	err := g.db.WithContext(ctx).Model(&BookCate{}).
		Where("deleted_at = ?", 0).
		Find(&bookCates).
		Error
	return bookCates, err
}

func (g *gormBookDao) InsertOrUpdate(ctx context.Context, id int64, name string) error {
	db := g.db.WithContext(ctx)
	now := time.Now().UnixMilli()
	if id == 0 {
		return db.Create(&BookCate{
			BaseModel: BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: 0,
			},
			Name: name,
		}).Error
	} else {
		return db.Model(&BookCate{}).Where("id = ?", id).Updates(map[string]any{
			"name":       name,
			"updated_at": now,
		}).Error
	}
}

func NewGormBookDao(db *gorm.DB, l logger.Logger) BookDao {
	return &gormBookDao{
		db: db,
		l:  l,
	}
}

type BookCate struct {
	BaseModel
	Name string
}

type Book struct {
	BaseModel
	BookName      string
	BookUrl       string
	BookAvatarUrl string
	BookCateId    int64
	BookUserId    int64
	BookMd5       string
	AvatarMd5     string
	Download      int64    `gorm:"default:0"`
	BookCate      BookCate `gorm:"foreignKey:BookCateId"`
}
