// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"github.com/WeiXinao/hi-wiki/internal/repository/cache"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/hi-wiki/internal/repository/oss"
	"github.com/WeiXinao/hi-wiki/internal/service"
	"github.com/WeiXinao/hi-wiki/internal/web"
	"github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() *App {
	handlerJWT := jwt.NewRedisJwtHandler()
	logger := InitLogger()
	v := InitMiddlewares(handlerJWT, logger)
	appConfig := InitViper()
	db := InitDB(appConfig, logger)
	userDao := dao.NewUserDao(db)
	prometheusHook := InitRedisPrometheus()
	otelHook := InitRedisOtel()
	cmdable := InitRedis(appConfig, prometheusHook, otelHook)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, userCache)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService, handlerJWT)
	fileDao := dao.NewFileDao(db)
	fileCache := cache.NewRedisFileCache(cmdable)
	fileRepository := repository.NewFileRepository(fileDao, fileCache)
	fileService := service.NewFileService(fileRepository)
	fileHandler := web.NewFileHandler(fileService)
	teamDao := dao.NewTeamDao(db)
	teamCache := cache.NewRedisTeamCache(cmdable)
	teamRepository := repository.NewTeamRepository(teamDao, teamCache)
	teamService := InitRetryableTeamService(teamRepository)
	teamHandler := web.NewTeamHandler(teamService)
	repoDao := dao.NewGormRepoDao(db)
	repoCache := cache.NewRedisRepoCache(cmdable)
	repoRepository := repository.NewRepoRepository(repoDao, repoCache)
	repoService := service.NewRepoService(repoRepository)
	repoHandler := web.NewRepoHandler(repoService)
	articleDao := dao.NewArticleDao(db)
	client := InitMinio(appConfig)
	ossOSS := oss.NewMinioOSS(client)
	articleRepository := repository.NewArticleRepository(articleDao, ossOSS, logger)
	articleService := service.NewArticleService(articleRepository, repoRepository)
	goodDao := dao.NewGoodDao(db)
	goodRepository := repository.NewGoodRepository(goodDao)
	goodService := service.NewGoodService(goodRepository)
	followDao := dao.NewGormFollowDao(db)
	followRepo := repository.NewFollowRepository(followDao)
	followService := service.NewFollowService(followRepo)
	articleHandler := web.NewArticleHandler(articleService, goodService, followService, logger)
	goodHandler := web.NewGoodHandler(goodService)
	followHandler := web.NewFollowHandler(followService)
	bookDao := dao.NewGormBookDao(db, logger)
	bookRepository := repository.NewBookRepository(bookDao)
	bookService := service.NewBookService(bookRepository)
	bookHandler := web.NewBookHandler(bookService, fileService)
	engine := InitServer(v, userHandler, fileHandler, teamHandler, repoHandler, articleHandler, goodHandler, followHandler, bookHandler)
	app := &App{
		Server: engine,
	}
	return app
}

// wire.go:

var thirdPartSet = wire.NewSet(
	InitViper,
	InitLogger,
	InitDB,
	InitRedis,
	InitRedisPrometheus,
	InitRedisOtel,
	InitMinio, jwt.NewRedisJwtHandler, InitMiddlewares, oss.NewMinioOSS,
)

var userHandlerSet = wire.NewSet(web.NewUserHandler, service.NewUserService, repository.NewUserRepository, cache.NewRedisUserCache, dao.NewUserDao)

var fileHandlerSet = wire.NewSet(web.NewFileHandler, service.NewFileService, repository.NewFileRepository, cache.NewRedisFileCache, dao.NewFileDao)

var teamHandlerSet = wire.NewSet(web.NewTeamHandler, InitRetryableTeamService, repository.NewTeamRepository, dao.NewTeamDao, cache.NewRedisTeamCache)

var repoHandlerSet = wire.NewSet(web.NewRepoHandler, service.NewRepoService, repository.NewRepoRepository, dao.NewGormRepoDao, cache.NewRedisRepoCache)

var articleHandlerSet = wire.NewSet(web.NewArticleHandler, service.NewArticleService, repository.NewArticleRepository, dao.NewArticleDao)

var goodHandlerSet = wire.NewSet(web.NewGoodHandler, service.NewGoodService, repository.NewGoodRepository, dao.NewGoodDao)

var followHandlerSet = wire.NewSet(web.NewFollowHandler, service.NewFollowService, repository.NewFollowRepository, dao.NewGormFollowDao)

var bookHandlerSet = wire.NewSet(web.NewBookHandler, service.NewBookService, repository.NewBookRepository, dao.NewGormBookDao)
