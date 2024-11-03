package vo

type ListRepoVO struct {
	Id         int64  `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	CateId     int64  `json:"cateid,omitempty"`
	CateName   string `json:"catename,omitempty"`
	UniqueCode string `json:"repo_unique_code,omitempty"`
	TeamId     int64  `json:"groupid,omitempty"`
	Desc       string `json:"repo_desc,omitempty"`
	CreateTime string `json:"createtime,omitempty"`
}
