package ioc

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/web"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/internal/web/middlewares"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	loggerMdl "github.com/WeiXinao/hi-wiki/pkg/ginx/middlewares/logger"
	"github.com/WeiXinao/hi-wiki/pkg/ginx/middlewares/prometheus"
	"github.com/WeiXinao/hi-wiki/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"strings"
)

func InitServer(
	middlewares []gin.HandlerFunc, userHandler *web.UserHandler,
	fileHandler *web.FileHandler, teamHandler *web.TeamHandler,
	repoHandler *web.RepoHandler, artsHandler *web.ArticleHandler,
	goodHandler *web.GoodHandler, followHandler *web.FollowHandler,
	bookHandler *web.BookHandler,
) *gin.Engine {
	server := gin.Default()

	server.Use(middlewares...)

	userHandler.RegisterRoutes(server)
	fileHandler.RegisterRoutes(server)
	teamHandler.RegisterRoutes(server)
	repoHandler.RegisterRoutes(server)
	artsHandler.RegisterRoutes(server)
	goodHandler.RegisterRoutes(server)
	followHandler.RegisterRoutes(server)
	bookHandler.RegisterRoutes(server)

	return server
}

func InitMiddlewares(
	jwtHdl _jwt.HandlerJWT,
	l logger.Logger,
) []gin.HandlerFunc {
	bd := loggerMdl.NewBuilder(func(ctx context.Context, al *loggerMdl.AccessLog) {
		l.Debug("HTTP",
			logger.String("method", al.Method),
			logger.String("url", al.Url),
			logger.Int("status", al.Status),
			logger.String("duration", al.Duration),
			logger.String("requestBody", al.ReqBody),
			logger.String("responseBody", al.RespBody),
		)
	}).AllowReqBody(true).AllowRespBody()

	pb := &prometheus.Builder{
		Namespace:  "xiaoxin",
		Subsystem:  "hi_wiki",
		Name:       "gin_http",
		Help:       "统计 GIN 的 HTTP 接口数据",
		InstanceId: "hi-wiki-1",
	}

	ginx.InitWrapperLogger(l)

	ginx.InitCounter(prom.CounterOpts{
		Namespace: "xiaoxin",
		Subsystem: "hi_wiki",
		Name:      "biz_code",
		Help:      "统计义务错误码",
	})

	return []gin.HandlerFunc{
		corsHandler(),
		bd.Build(),
		middlewares.NewAuthenticationBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/files/img/:md5").
			IgnorePaths("/books/download/:id").
			Build(),
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		otelgin.Middleware("hi_wiki"),
	}
}

func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
	})
}
