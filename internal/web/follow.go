package web

import (
	"errors"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/internal/service"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	"github.com/gin-gonic/gin"
	"strings"
)

type FollowHandler struct {
	svc service.FollowService
}

func NewFollowHandler(svc service.FollowService) *FollowHandler {
	return &FollowHandler{
		svc: svc,
	}
}

func (h *FollowHandler) RegisterRoutes(server *gin.Engine) {
	fg := server.Group("/follows")
	fg.POST("/add", ginx.WrapBodyAndClaims[FollowReq, *_jwt.UserClaims](h.Add))
	fg.GET("/recent", ginx.WrapClaims[*_jwt.UserClaims](h.Recent))
}

type FollowReq struct {
	Flag string `json:"flag"`
	Id   int64  `json:"id"`
}

func (h *FollowHandler) Recent(ctx *gin.Context, uc *_jwt.UserClaims) (ginx.Result, error) {
	recent, err := h.svc.Recent(ctx, uc.UserId)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessShowRecentFollow, recent), nil
}

func (h *FollowHandler) Add(ctx *gin.Context, req FollowReq,
	uc *_jwt.UserClaims) (ginx.Result, error) {
	typ := domain.FollowTypeFactory(strings.TrimSpace(req.Flag))
	if !typ.Valid() {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	err := h.svc.Add(ctx.Request.Context(), typ, req.Id, uc.UserId)
	if errors.Is(err, service.ErrFollowedObjectNotExists) {
		return MarshalResp(errs.FollowedObjectNotExists, nil), nil
	}
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessFollow, nil), nil
}
