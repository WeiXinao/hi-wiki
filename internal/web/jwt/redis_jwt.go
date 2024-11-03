package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

var ErrIllegalAuthorizationHeader = errors.New("非法Authorization头")

func NewRedisJwtHandler() HandlerJWT {
	return &redisJwtHandler{}
}

type redisJwtHandler struct {
}

func (r *redisJwtHandler) GetJwtToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	headerSegments := strings.Split(header, " ")
	if len(headerSegments) != 2 {
		return "", ErrIllegalAuthorizationHeader
	}
	return headerSegments[1], nil
}

func (r *redisJwtHandler) SetJwtToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 一个小时不续期，jwt-token 就过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UserId: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", signedToken)
	return nil
}
