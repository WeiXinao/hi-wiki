package domain

import "time"

type File struct {
	Id             int64
	Name           string
	Typ            string
	Md5            string
	Url            string
	Size           string
	DirLevel       uint
	RepoUniqueCode string
	Owner          Owner
	CreateAt       time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type Owner struct {
	Id       int64
	Username string
}
