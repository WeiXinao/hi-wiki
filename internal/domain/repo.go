package domain

import (
	"github.com/WeiXinao/xkit/slice"
	"time"
)

type Repo struct {
	Id         int64
	Name       string
	Desc       string
	Cate       Cate
	Status     RepoStatus
	UniqueCode string
	State      uint8
	UserGroup  UserGroup
	Creator    Creator
	CreateAt   time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

type RepoStatus uint8

const (
	RepoStatusOnlyMemberAvailable RepoStatus = iota + 1
	RepoStatusAllAvailable
)

var RepoStatusList = []RepoStatus{RepoStatusOnlyMemberAvailable, RepoStatusAllAvailable}

func (t RepoStatus) ToUint8() uint8 {
	return uint8(t)
}

func (t RepoStatus) Valid() bool {
	return slice.Contains[RepoStatus](RepoStatusList, t)
}

func (t RepoStatus) IsPrivate() bool {
	return t == RepoStatusOnlyMemberAvailable
}

type UserGroup struct {
	Id         int64
	Name       string
	Desc       string
	UniqueCode string
	Status     RepoStatus
	AvatarMd5  string
}

type Creator struct {
	Id        int64
	Username  string
	AvatarMd5 string
	Profile   string
}

type Cate struct {
	Id   int64
	Name string
}
