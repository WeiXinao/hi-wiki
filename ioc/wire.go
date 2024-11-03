//go:build wireinject

package ioc

import (
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"github.com/WeiXinao/hi-wiki/internal/repository/cache"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/hi-wiki/internal/repository/oss"
	"github.com/WeiXinao/hi-wiki/internal/service"
	"github.com/WeiXinao/hi-wiki/internal/web"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/google/wire"
)

var thirdPartSet = wire.NewSet(
	InitViper,
	InitLogger,
	InitDB,
	InitRedis,
	InitRedisPrometheus,
	InitRedisOtel,
	InitMinio,
	_jwt.NewRedisJwtHandler,
	InitMiddlewares,
	oss.NewMinioOSS,
)

var userHandlerSet = wire.NewSet(
	web.NewUserHandler,
	service.NewUserService,
	repository.NewUserRepository,
	cache.NewRedisUserCache,
	dao.NewUserDao,
)

var fileHandlerSet = wire.NewSet(
	web.NewFileHandler,
	service.NewFileService,
	repository.NewFileRepository,
	cache.NewRedisFileCache,
	dao.NewFileDao,
)

var teamHandlerSet = wire.NewSet(
	web.NewTeamHandler,
	InitRetryableTeamService,
	repository.NewTeamRepository,
	dao.NewTeamDao,
	cache.NewRedisTeamCache,
)

var repoHandlerSet = wire.NewSet(
	web.NewRepoHandler,
	service.NewRepoService,
	repository.NewRepoRepository,
	dao.NewGormRepoDao,
	cache.NewRedisRepoCache,
)

var articleHandlerSet = wire.NewSet(
	web.NewArticleHandler,
	service.NewArticleService,
	repository.NewArticleRepository,
	dao.NewArticleDao,
)

var goodHandlerSet = wire.NewSet(
	web.NewGoodHandler,
	service.NewGoodService,
	repository.NewGoodRepository,
	dao.NewGoodDao,
)

var followHandlerSet = wire.NewSet(
	web.NewFollowHandler,
	service.NewFollowService,
	repository.NewFollowRepository,
	dao.NewGormFollowDao,
)

var bookHandlerSet = wire.NewSet(
	web.NewBookHandler,
	service.NewBookService,
	repository.NewBookRepository,
	dao.NewGormBookDao,
)

func InitApp() *App {
	wire.Build(
		// 第三方依赖
		thirdPartSet,
		userHandlerSet,
		fileHandlerSet,
		teamHandlerSet,
		repoHandlerSet,
		articleHandlerSet,
		goodHandlerSet,
		followHandlerSet,
		bookHandlerSet,

		InitServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
