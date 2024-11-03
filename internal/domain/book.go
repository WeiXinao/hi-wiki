package domain

import "time"

type BookCate struct {
	Id        int64
	Name      string
	CreateAt  time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Book struct {
	Id            int64
	BookName      string
	BookUrl       string
	BookAvatarUrl string
	BookCateId    int64
	BookCateName  string
	BookUserId    int64
	BookMd5       string
	AvatarMd5     string
	Download      int64
	CreateAt      time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

func (b *Book) GetAvatarUrl() string {
	return "files/img/" + b.AvatarMd5
}
