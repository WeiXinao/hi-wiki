package vo

type UserVO struct {
	Id        int64  `json:"userid"`
	Username  string `json:"username"`
	AvatarMd5 string `json:"avatarMd5"`
	AvatarUrl string `json:"avatarurl"`
	Profile   string `json:"desc"`
}
