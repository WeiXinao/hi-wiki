package ginx

import (
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

var (
	wrapLog logger.Logger
	vector  *prometheus.CounterVec
)

func InitWrapperLogger(l logger.Logger) {
	wrapLog = l
}

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}

func Wrap(fn func(ctx *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			// 记录日志
			wrapLog.Error("处理业务逻辑出错",
				// 请求的路径
				logger.String("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBody[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}
		res, err := fn(ctx, req)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			wrapLog.Error("处理业务逻辑出错",
				logger.String("path", ctx.Request.URL.Path),
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapClaims[C jwt.Claims](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		userClaims, ok := claims.(C)
		if !exists || !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		res, err := fn(ctx, userClaims)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			// 记录日志
			wrapLog.Error("处理业务逻辑出错",
				// 请求的路径
				logger.String("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBodyAndClaims[T any, C jwt.Claims](fn func(ctx *gin.Context, req T, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}

		claims, exists := ctx.Get("claims")
		userClaims, ok := claims.(C)
		if !exists || !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		res, err := fn(ctx, req, userClaims)
		vector.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		if err != nil {
			// 记录日志
			wrapLog.Error("处理业务逻辑出错",
				// 请求的路径
				logger.String("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

type Result struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

func MarshalResp(err errs.Err, data any) Result {
	return Result{
		Code: err.GetCode(),
		Msg:  err.GetMsg(),
		Data: data,
	}
}
