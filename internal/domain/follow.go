package domain

import (
	"github.com/WeiXinao/xkit/slice"
	"strings"
	"time"
)

type FollowInfo struct {
	FollowCount   int64
	BeFollowCount int64
	IsFollowed    bool
}

type Follow struct {
	Id         int64
	FollowType uint8
	UserId     int64
	FollowId   int64
	CreateAt   time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

type FollowRecent struct {
	UserName   string `json:"username,omitempty"`
	Avatar     string `json:"avator,omitempty"`
	UniqueCode string `json:"unique_code,omitempty"`
	Title      string `json:"title,omitempty"`
	CreateTime string `json:"createtime,omitempty"`
	Typ        string `json:"flag,omitempty"`
}

type FollowType uint8

const (
	FollowTypeUser FollowType = iota
	FollowTypeRepo
	FollowTypeInvalid
)

var FollowTypeList = []FollowType{FollowTypeUser, FollowTypeRepo}

func FollowTypeFactory(typ string) FollowType {
	switch strings.TrimSpace(typ) {
	case "user":
		return FollowTypeUser
	case "repo":
		return FollowTypeRepo
	default:
		return FollowTypeInvalid
	}
}

func (f FollowType) String() string {
	switch f {
	case FollowTypeUser:
		return "user"
	case FollowTypeRepo:
		return "repo"
	default:
		return ""
	}
}

func (f FollowType) ToUint8() uint8 {
	return uint8(f)
}

func (f FollowType) Valid() bool {
	return slice.Contains[FollowType](FollowTypeList, f)
}
