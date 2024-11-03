package vo

import "time"

type TeamVO struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	UniqueCode string `json:"flag"`
	Desc       string `json:"groupdesc"`
	AvatarMd5  string `json:"avatorurl"`
}

// TeamInfoVO 获取团队的详细信息
type TeamInfoVO struct {
	Id         int64     `json:"id"`
	Name       string    `json:"group_name"`
	UniqueCode string    `json:"group_unique_code"`
	Desc       string    `json:"group_desc"`
	Status     uint8     `json:"group_status"`
	AvatarMd5  string    `json:"group_avator_url_id"`
	CreatedAt  time.Time `json:"createtime"`
}

type TeamMemberVO struct {
	Uid      int64  `json:"uid,omitempty"`
	Username string `json:"username,omitempty"`
	UserType string `json:"usertype,omitempty"`
}
