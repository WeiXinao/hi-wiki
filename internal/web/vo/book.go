package vo

type BookCateVO struct {
	Id       int64  `json:"id,omitempty"`
	CateName string `json:"catename,omitempty"`
}

type BookVO struct {
	Id            int64  `json:"id"`
	BookName      string `json:"name"`
	BookUrl       string `json:"fileurl"`
	BookAvatarUrl string `json:"avatorurl"`
	BookCateName  string `json:"catename"`
	Download      int64  `json:"download"`
}
