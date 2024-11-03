package domain

import (
	"github.com/WeiXinao/xkit/slice"
	"time"
)

type Team struct {
	Id         int64
	Name       string
	Desc       string
	UniqueCode string
	Status     TeamStatus
	AvatarMd5  string
	CreateAt   time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

type TeamMember struct {
	Id        int64
	Name      string
	IsLeader  bool
	UserId    int64
	CreateAt  time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type TeamStatus uint8

const (
	TeamStatusOnlyMemberAvailable TeamStatus = iota + 1
	TeamStatusAllAvailable
)

var TeamStatusList = []TeamStatus{TeamStatusOnlyMemberAvailable, TeamStatusAllAvailable}

func (t TeamStatus) ToUint8() uint8 {
	return uint8(t)
}

func (t TeamStatus) Valid() bool {
	return slice.Contains[TeamStatus](TeamStatusList, t)
}
