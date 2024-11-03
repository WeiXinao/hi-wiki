package middlewares

import (
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/xkit/slice"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type AuthenticationBuilder struct {
	paths []string
	_jwt.HandlerJWT
}

func NewAuthenticationBuilder(jwtHdl _jwt.HandlerJWT) *AuthenticationBuilder {
	return &AuthenticationBuilder{
		HandlerJWT: jwtHdl,
	}
}

func (a *AuthenticationBuilder) IgnorePaths(path string) *AuthenticationBuilder {
	a.paths = append(a.paths, path)
	return a
}

func (a *AuthenticationBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//	检验白名单
		if slice.Contains[string](a.paths, ctx.FullPath()) {
			return
		}

		authorizationHeader, err := a.GetJwtToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &_jwt.UserClaims{}
		token, err := jwt.ParseWithClaims(authorizationHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return _jwt.AtKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		if token == nil || !token.Valid || claims.UserId == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 续期
		now := time.Now()
		// 每十分钟续期一次
		if claims.ExpiresAt.Sub(now) < time.Minute*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(10 * time.Minute))
			authorizationHeader, err = token.SignedString(_jwt.AtKey)
			if err != nil {
				//	记录日志
			}
			ctx.Header("x-jwt-token", authorizationHeader)
		}

		ctx.Set("claims", claims)
	}
}
