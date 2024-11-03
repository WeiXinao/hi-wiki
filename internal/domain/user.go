package domain

import "time"

type User struct {
	Id        int64
	Username  string
	Password  string
	AvatarMd5 string
	Profile   string
	CreateAt  time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
