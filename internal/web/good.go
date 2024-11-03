package web

import (
	"errors"
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/internal/service"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	"github.com/gin-gonic/gin"
	"strings"
)

type GoodHandler struct {
	svc service.GoodService
}

func NewGoodHandler(svc service.GoodService) *GoodHandler {
	return &GoodHandler{
		svc: svc,
	}
}

func (h *GoodHandler) RegisterRoutes(server *gin.Engine) {
	gg := server.Group("goods")
	gg.POST("/article", ginx.WrapBodyAndClaims[GoodReq, *_jwt.UserClaims](h.GoodArticle))
}

type GoodReq struct {
	ArticleUniqueCode string `json:"uniquecode"`
}

func (h *GoodHandler) GoodArticle(ctx *gin.Context, req GoodReq, uc *_jwt.UserClaims) (ginx.Result, error) {
	uniqueCode := req.ArticleUniqueCode
	if len(strings.TrimSpace(uniqueCode)) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	err := h.svc.GoodArticle(ctx.Request.Context(), uniqueCode, uc.UserId)
	if errors.Is(err, service.ErrHasGood) {
		return MarshalResp(errs.HasGood, nil), nil
	}
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessGoodArticle, nil), nil
}
