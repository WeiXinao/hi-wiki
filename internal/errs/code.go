package errs

// 通用模块统一错误码
var (
	InternalServerError = bizErr{
		code: 500001,
		msg:  "系统错误",
	}
	InternalInvalidInput = bizErr{
		code: 400001,
		msg:  "无效输入",
	}
)

// 用户模块的统一错误码
var (
	SuccessLogin = bizErr{
		code: 201001,
		msg:  "登录成功",
	}
	SuccessSign = bizErr{
		code: 201002,
		msg:  "注册成功",
	}
	SuccessModifyPassword = bizErr{
		code: 201003,
		msg:  "修改密码成功",
	}
	SuccessModifyProfile = bizErr{
		code: 201004,
		msg:  "修改个人信息成功",
	}
	SuccessGetProfile = bizErr{
		code: 201005,
		msg:  "获取个人信息成功",
	}

	InvalidUserOrPassword = bizErr{
		code: 401001,
		msg:  "用户名或密码错误",
	}
	UserDuplicatedUsername = bizErr{
		code: 401002,
		msg:  "用户名已存在",
	}
	InconsistentPasswordAndConfirmPassword = bizErr{
		code: 401003,
		msg:  "密码和确认密码不一致",
	}
	UsernameOutOfRange = bizErr{
		code: 401004,
		msg:  "用户名长度范围不对",
	}
	InconsistentTwoPassword = bizErr{
		code: 401005,
		msg:  "两次密码不一致",
	}
	InvalidPassword = bizErr{
		code: 401006,
		msg:  "密码错误",
	}
	FailModifyProfile = bizErr{
		code: 501001,
		msg:  "修改个人信息失败",
	}
)

// 文件模块的统一错误码
var (
	SuccessUploadFile = bizErr{
		code: 202001,
		msg:  "上传文件成功",
	}
	ImageNotFound = bizErr{
		code: 402001,
		msg:  "图片不存在",
	}
	FailUploadFile = bizErr{
		code: 502001,
		msg:  "上传文件失败",
	}
)

// 团队模块的统一错误码
var (
	SuccessCreateTeam = bizErr{
		code: 203001,
		msg:  "创建团队成功",
	}
	SuccessListTeam = bizErr{
		code: 203002,
		msg:  "获取团队列表成功",
	}
	SuccessShowDetail = bizErr{
		code: 203003,
		msg:  "展示团队细节成功",
	}
	SuccessListTeamMember = bizErr{
		code: 203004,
		msg:  "获取团队成员成功",
	}
	SuccessDeleteTeamMember = bizErr{
		code: 203005,
		msg:  "删除成员成功",
	}
	TeamLeaderCannotBeDeleted = bizErr{
		code: 503001,
		msg:  "不能删除组长",
	}
)

// 知识库模块的统一错误码
var (
	SuccessCreateRepo = bizErr{
		code: 204001,
		msg:  "创建知识库成功",
	}
	SuccessListRepo = bizErr{
		code: 204002,
		msg:  "展示知识库列表成功",
	}
)

// 文章模块统一错误码
var (
	SuccessEditArticle = bizErr{
		code: 205001,
		msg:  "编辑文章成功",
	}
	SuccessListArticle = bizErr{
		code: 205002,
		msg:  "展示文章列表成功",
	}
	SuccessShowArticleDetail = bizErr{
		code: 205003,
		msg:  "展示文章详情成功",
	}
)

// 点赞模块统一错误码
var (
	SuccessGoodArticle = bizErr{
		code: 206001,
		msg:  "点赞成功",
	}
	HasGood = bizErr{
		code: 406001,
		msg:  "已经点过赞了",
	}
)

// 关注模块的统一错误码
var (
	SuccessFollow = bizErr{
		code: 207001,
		msg:  "关注成功",
	}
	SuccessShowRecentFollow = bizErr{
		code: 207002,
		msg:  "展示关注动态成功",
	}
	FollowedObjectNotExists = bizErr{
		code: 407001,
		msg:  "关注的对象不存在",
	}
)

// 图书模块的统一错误码
var (
	SuccessEditBookCate = bizErr{
		code: 208001,
		msg:  "成功编辑图书分类",
	}
	SuccessListBookCates = bizErr{
		code: 208002,
		msg:  "成功展示数据分类",
	}
	SuccessDeleteCate = bizErr{
		code: 208003,
		msg:  "成功删除分类",
	}
	SuccessUploadBook = bizErr{
		code: 208004,
		msg:  "成功上传图书",
	}
	SuccessListBook = bizErr{
		code: 208005,
		msg:  "成功展示书籍",
	}
	DeletedBookCateNotFound = bizErr{
		code: 408001,
		msg:  "待删除的分类未找到",
	}
	FailAddBookInfo = bizErr{
		code: 508001,
		msg:  "添加图书信息失败",
	}
)

type Err interface {
	GetCode() int
	GetMsg() string
}

type bizErr struct {
	code int
	msg  string
}

func (be bizErr) GetCode() int {
	return be.code
}

func (be bizErr) GetMsg() string {
	return be.msg
}
