package domain

import "time"

type GoodHistory struct {
	Id                int64
	UserId            int64
	ArticleUniqueCode string `gorm:"varchar(191)"`
	CreateAt          time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}
