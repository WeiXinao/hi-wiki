package service

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
)

var ErrDeletedBookCateNotFound = repository.ErrDeletedBookCateNotFound

type BookService interface {
	EditCate(ctx context.Context, id int64, name string) error
	ListCate(ctx context.Context) ([]domain.BookCate, error)
	DelCate(ctx context.Context, id int64) error
	Add(ctx context.Context, bookName string, bookUrl string, bookAvatarUrl string,
		bookCateId int64, bookUserId int64, bookMd5 string, avatarMd5 string) error
	List(ctx context.Context, cate int64, kw string, rank string, page int, offset int) ([]domain.Book, error)
	Download(ctx context.Context, id int64) (domain.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

func (b *bookService) Download(ctx context.Context, id int64) (domain.Book, error) {
	book, err := b.repo.UpdateAndGet(ctx, id)
	if err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func (b *bookService) List(ctx context.Context, cate int64, kw string, rank string, page int, size int) ([]domain.Book, error) {
	return b.repo.Get(ctx, cate, kw, rank, page, size)
}

func (b *bookService) Add(ctx context.Context, bookName string, bookUrl string, bookAvatarUrl string,
	bookCateId int64, bookUserId int64, bookMd5 string, avatarMd5 string) error {
	return b.repo.InsertCate(ctx, bookName, bookUrl, bookAvatarUrl, bookCateId, bookUserId, bookMd5, avatarMd5)
}

func (b *bookService) DelCate(ctx context.Context, id int64) error {
	return b.repo.DelCate(ctx, id)
}

func (b *bookService) ListCate(ctx context.Context) ([]domain.BookCate, error) {
	return b.repo.GetCate(ctx)
}

func (b *bookService) EditCate(ctx context.Context, id int64, name string) error {
	return b.repo.InsertOrUpdateCate(ctx, id, name)
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{
		repo: repo,
	}
}
