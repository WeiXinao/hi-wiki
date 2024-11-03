package repository

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/xkit/slice"
	"time"
)

var ErrDeletedBookCateNotFound = dao.ErrDeletedBookCateNotFound

type BookRepository interface {
	InsertOrUpdateCate(ctx context.Context, id int64, name string) error
	GetCate(ctx context.Context) ([]domain.BookCate, error)
	DelCate(ctx context.Context, id int64) error
	InsertCate(ctx context.Context, bookName string, bookUrl string, bookAvatarUrl string,
		bookCateId int64, bookUserId int64, bookMd5 string, avatarMd5 string) error
	Get(ctx context.Context, cate int64, kw string, rank string, page int, size int) ([]domain.Book, error)
	UpdateAndGet(ctx context.Context, id int64) (domain.Book, error)
}

type bookRepository struct {
	dao dao.BookDao
}

func (b *bookRepository) UpdateAndGet(ctx context.Context, id int64) (domain.Book, error) {
	book, err := b.dao.UpdateAndGet(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}
	return b.BookEntityToDomain(book), err
}

func (b *bookRepository) Get(ctx context.Context, cate int64, kw string, rank string, page int, size int) ([]domain.Book, error) {
	books, err := b.dao.Get(ctx, cate, kw, rank, page, size)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Book, domain.Book](books, func(idx int, src dao.Book) domain.Book {
		return b.BookEntityToDomain(src)
	}), err
}

func (b *bookRepository) InsertCate(ctx context.Context, bookName string, bookUrl string, bookAvatarUrl string,
	bookCateId int64, bookUserId int64, bookMd5 string, avatarMd5 string) error {
	return b.dao.Insert(ctx, bookName, bookUrl, bookAvatarUrl,
		bookCateId, bookUserId, bookMd5, avatarMd5)
}

func (b *bookRepository) DelCate(ctx context.Context, id int64) error {
	return b.dao.Del(ctx, id)
}

func (b *bookRepository) GetCate(ctx context.Context) ([]domain.BookCate, error) {
	bookCates, err := b.dao.List(ctx)
	return slice.Map[dao.BookCate, domain.BookCate](bookCates,
		func(idx int, src dao.BookCate) domain.BookCate {
			return b.BookCateEntityToDomain(src)
		}), err
}

func (b *bookRepository) InsertOrUpdateCate(ctx context.Context, id int64, name string) error {
	return b.dao.InsertOrUpdate(ctx, id, name)
}

func (b *bookRepository) BookEntityToDomain(book dao.Book) domain.Book {
	return domain.Book{
		Id:            book.Id,
		BookName:      book.BookName,
		BookUrl:       book.BookUrl,
		BookAvatarUrl: book.BookAvatarUrl,
		BookCateId:    book.BookCateId,
		BookCateName:  book.BookCate.Name,
		BookUserId:    book.BookUserId,
		BookMd5:       book.BookMd5,
		AvatarMd5:     book.AvatarMd5,
		Download:      book.Download,
		CreateAt:      time.UnixMilli(book.CreatedAt),
		UpdatedAt:     time.UnixMilli(book.UpdatedAt),
		DeletedAt:     time.UnixMilli(book.DeletedAt),
	}
}

func (b *bookRepository) BookCateEntityToDomain(bookCate dao.BookCate) domain.BookCate {
	return domain.BookCate{
		Id:        bookCate.Id,
		Name:      bookCate.Name,
		CreateAt:  time.UnixMilli(bookCate.CreatedAt),
		UpdatedAt: time.UnixMilli(bookCate.UpdatedAt),
		DeletedAt: time.UnixMilli(bookCate.DeletedAt),
	}
}

func NewBookRepository(dao dao.BookDao) BookRepository {
	return &bookRepository{
		dao: dao,
	}
}
