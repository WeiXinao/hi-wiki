package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type HandlerJWT interface {
	SetJwtToken(ctx *gin.Context, uid int64) error
	GetJwtToken(ctx *gin.Context) (string, error)
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64
}
