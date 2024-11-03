package vo

type ArticleVO struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	LikeCnt    int64  `json:"liked"`
	UniqueCode string `json:"ar_unique_code"`
	UserId     int64  `json:"userid"`
	IsFollow   bool   `json:"isFollow"`
	IsGood     bool   `json:"isGood"`
	Cate       string `json:"catename"`
	CreateAt   string `json:"createtime"`
	User       User   `json:"user"`
}

type ArticleWithRepoVO struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	RepoName string `json:"reponame"`
	RepoCode string `json:"repocode"`
	CreateAt string `json:"createtime"`
	IsAuthor bool   `json:"isauthor"`
}

type User struct {
	Username      string `json:"username"`
	AvatarMd5     string `json:"avatarurl"`
	Profile       string `json:"profile"`
	FollowCount   int64  `json:"followCount"`
	BeFollowCount int64  `json:"beFollowCount"`
}
