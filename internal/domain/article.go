package domain

import "time"

type Article struct {
	Id             int64
	Title          string
	Content        string
	PureContent    string
	Desc           string
	LikeCnt        int64
	UniqueCode     string
	RepoUniqueCode string
	Author         Author
	Cate           Cate
	Repo           Repo
	State          uint8
	Private        bool
	CreateAt       time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

func (a *Article) GetDesc() string {
	return a.Desc
}

func (a *Article) GenDesc() {
	contentLen := len([]rune(a.PureContent))
	if contentLen < 50 {
		a.Desc = a.PureContent
		if contentLen == 0 {
			a.Desc = "文章暂无介绍信息"
		}
	} else {
		a.Desc = string([]rune(a.PureContent)[:50])
	}
}

func (a *Article) GetFormatCreateTime() string {
	return a.CreateAt.Format("2006-01-02 15:04")
}

type Author struct {
	Id        int64
	Name      string
	AvatarMd5 string
	Profile   string
}

func (a Author) GetAvatarUrl() string {
	return "files/img/" + a.AvatarMd5
}
