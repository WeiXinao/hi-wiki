package web

import (
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/internal/service"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/internal/web/vo"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	"github.com/WeiXinao/hi-wiki/pkg/logger"
	"github.com/WeiXinao/xkit/slice"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type ArticleHandler struct {
	svc       service.ArticleService
	goodSvc   service.GoodService
	followSvc service.FollowService
	l         logger.Logger
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	ag := server.Group("/arts")
	ag.POST("/edit", ginx.WrapBodyAndClaims[ArticleReq, *_jwt.UserClaims](h.Edit))
	ag.GET("/list/self", ginx.WrapClaims[*_jwt.UserClaims](h.ListSelf))
	ag.GET("/detail/:id", ginx.WrapClaims[*_jwt.UserClaims](h.Detail))
	ag.GET("/list/all", ginx.WrapClaims[*_jwt.UserClaims](h.ListAll))
}

func (h *ArticleHandler) ListAll(ctx *gin.Context, uc *_jwt.UserClaims) (ginx.Result, error) {
	offset, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	hot := ctx.Query("flag")
	arts, count, err := h.svc.ListAll(ctx.Request.Context(), offset, size,
		strings.TrimSpace(hot) == "hot")
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	ids := slice.Map[domain.Article, int64](arts, func(idx int, src domain.Article) int64 {
		return src.Author.Id
	})
	followType := domain.FollowTypeFactory("user")
	info, err := h.followSvc.StatsFollowInfo(ctx.Request.Context(), ids, uc.UserId, ids, followType.ToUint8())
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	likedMap, err := h.goodSvc.IsLiked(ctx.Request.Context(), uc.UserId,
		slice.Map[domain.Article, string](arts, func(idx int, src domain.Article) string {
			return src.UniqueCode
		}))
	if err != nil {
		h.l.Error("判断是否点赞失败", logger.Error(err))
	}
	data := map[string]any{
		"count": count,
		"arts": slice.Map[domain.Article, vo.ArticleVO](arts, func(idx int, src domain.Article) vo.ArticleVO {

			return vo.ArticleVO{
				Id:         src.Id,
				Title:      src.Title,
				Desc:       src.Desc,
				LikeCnt:    src.LikeCnt,
				UniqueCode: src.UniqueCode,
				UserId:     src.Author.Id,
				IsFollow:   info[src.Author.Id].IsFollowed,
				IsGood:     likedMap[src.UniqueCode],
				Cate:       src.Cate.Name,
				CreateAt:   src.GetFormatCreateTime(),
				User: vo.User{
					Username:      src.Author.Name,
					AvatarMd5:     src.Author.GetAvatarUrl(),
					Profile:       src.Author.Profile,
					FollowCount:   info[src.Author.Id].FollowCount,
					BeFollowCount: info[src.Author.Id].BeFollowCount,
				},
			}
		}),
	}
	return MarshalResp(errs.SuccessListArticle, data), nil
}

func (h *ArticleHandler) Detail(ctx *gin.Context, uc *_jwt.UserClaims) (ginx.Result, error) {
	uniqueCode := ctx.Param("id")
	if len(strings.TrimSpace(uniqueCode)) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	art, err := h.svc.Detail(ctx, uniqueCode, uc.UserId)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	artVO := vo.ArticleWithRepoVO{
		Title:    art.Title,
		Content:  art.Content,
		RepoName: art.Repo.Name,
		RepoCode: art.Repo.UniqueCode,
		CreateAt: art.GetFormatCreateTime(),
		IsAuthor: art.Author.Id == uc.UserId,
	}
	return MarshalResp(errs.SuccessShowArticleDetail, artVO), nil
}

func (h *ArticleHandler) ListSelf(ctx *gin.Context, uc *_jwt.UserClaims) (ginx.Result, error) {
	offset, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	arts, count, err := h.svc.ListSelf(ctx.Request.Context(), uc.UserId, offset, size)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	data := map[string]any{
		"count": count,
		"arts": slice.Map[domain.Article, vo.ArticleVO](arts, func(idx int, src domain.Article) vo.ArticleVO {
			return vo.ArticleVO{
				Id:         src.Id,
				Title:      src.Title,
				Desc:       src.Desc,
				LikeCnt:    src.LikeCnt,
				UniqueCode: src.UniqueCode,
				UserId:     src.Author.Id,
				Cate:       src.Cate.Name,
				CreateAt:   src.GetFormatCreateTime(),
				User: vo.User{
					Username:  src.Author.Name,
					AvatarMd5: src.Author.GetAvatarUrl(),
					Profile:   src.Author.Profile,
				},
			}
		}),
	}
	return MarshalResp(errs.SuccessListArticle, data), nil
}

type ArticleReq struct {
	Title          string `json:"title,omitempty"`
	RepoUniqueCode string `json:"repoid,omitempty"`
	Content        string `json:"content,omitempty"`
	PureText       string `json:"puretext,omitempty"`
	UniqueCode     string `json:"aid,omitempty"`
}

func (h *ArticleHandler) Edit(ctx *gin.Context, req ArticleReq, uc *_jwt.UserClaims) (ginx.Result, error) {
	article := domain.Article{
		Title:          req.Title,
		Content:        req.Content,
		PureContent:    req.PureText,
		UniqueCode:     req.UniqueCode,
		RepoUniqueCode: req.RepoUniqueCode,
		Author: domain.Author{
			Id: uc.UserId,
		},
	}
	err := h.svc.Edit(ctx.Request.Context(), article)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessEditArticle, nil), nil
}

func NewArticleHandler(svc service.ArticleService,
	goodSvc service.GoodService,
	followSvc service.FollowService,
	l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:       svc,
		goodSvc:   goodSvc,
		followSvc: followSvc,
		l:         l,
	}
}
