package web

import (
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/internal/service"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/internal/web/vo"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	"github.com/WeiXinao/xkit/slice"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type RepoHandler struct {
	svc service.RepoService
}

func NewRepoHandler(svc service.RepoService) *RepoHandler {
	return &RepoHandler{
		svc: svc,
	}
}

func (h *RepoHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/repos")
	rg.POST("/create", ginx.WrapBodyAndClaims[CreateRepoReq, *_jwt.UserClaims](h.Create))
	rg.GET("/list", ginx.WrapClaims[*_jwt.UserClaims](h.List))
	rg.GET("/list_by_team", ginx.WrapClaims[*_jwt.UserClaims](h.ListByTeamId))
	rg.GET("/list_hot", ginx.Wrap(h.ListHot))
}

type CreateRepoReq struct {
	Auth     uint8  `json:"auth"`
	Desc     string `json:"desc"`
	Reponame string `json:"reponame"`
	Repotype int64  `json:"repotype"`
	Team     int64  `json:"team"`
}

func (h *RepoHandler) ListHot(ctx *gin.Context) (ginx.Result, error) {
	var (
		repos []domain.Repo
		err   error
	)
	repoType := ctx.Query("type")
	flag := ctx.Query("flag")
	isDoc := strings.TrimSpace(repoType) == "doc"
	isHot := strings.TrimSpace(flag) == "hot"
	repos, err = h.svc.ListHot(ctx.Request.Context(), isDoc, isHot)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessListRepo, slice.Map[domain.Repo, vo.ListRepoVO](repos, func(idx int, src domain.Repo) vo.ListRepoVO {
		return vo.ListRepoVO{
			Id:         src.Id,
			Name:       src.Name,
			CateId:     src.Cate.Id,
			CateName:   src.Cate.Name,
			UniqueCode: src.UniqueCode,
			TeamId:     src.UserGroup.Id,
			Desc:       src.Desc,
			CreateTime: src.CreateAt.Format(time.DateOnly),
		}
	})), nil
}

func (h *RepoHandler) ListByTeamId(ctx *gin.Context, uc *_jwt.UserClaims) (Result, error) {
	var (
		repos []domain.Repo
		err   error
	)
	groupId := ctx.Query("groupid")
	IntGroupId, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	repos, err = h.svc.ListByTeamId(ctx.Request.Context(), IntGroupId)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessListRepo, slice.Map[domain.Repo, vo.ListRepoVO](repos, func(idx int, src domain.Repo) vo.ListRepoVO {
		return vo.ListRepoVO{
			Id:         src.Id,
			Name:       src.Name,
			CateId:     src.Cate.Id,
			CateName:   src.Cate.Name,
			UniqueCode: src.UniqueCode,
			TeamId:     src.UserGroup.Id,
			Desc:       src.Desc,
			CreateTime: src.CreateAt.Format(time.DateOnly),
		}
	})), nil
}

func (h *RepoHandler) List(ctx *gin.Context, uc *_jwt.UserClaims) (Result, error) {
	var (
		repos []domain.Repo
		err   error
	)
	repoType := ctx.Query("type")
	repos, err = h.svc.List(ctx.Request.Context(), uc.UserId, repoType == "doc")
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessListRepo, slice.Map[domain.Repo, vo.ListRepoVO](repos, func(idx int, src domain.Repo) vo.ListRepoVO {
		return vo.ListRepoVO{
			Id:         src.Id,
			Name:       src.Name,
			CateId:     src.Cate.Id,
			CateName:   src.Cate.Name,
			UniqueCode: src.UniqueCode,
			TeamId:     src.UserGroup.Id,
			Desc:       src.Desc,
			CreateTime: src.CreateAt.Format(time.DateOnly),
		}
	})), nil
}

func (h *RepoHandler) Create(ctx *gin.Context, req CreateRepoReq, uc *_jwt.UserClaims) (Result, error) {
	// 参数校验
	status := domain.RepoStatus(req.Auth)
	if !status.Valid() {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	err := h.svc.Create(ctx.Request.Context(), domain.Repo{
		Name: req.Reponame,
		Desc: req.Desc,
		Cate: domain.Cate{
			Id: req.Repotype,
		},
		Status:     status,
		UniqueCode: "",
		State:      0,
		UserGroup: domain.UserGroup{
			Id: req.Team,
		},
		Creator: domain.Creator{
			Id: uc.UserId,
		},
	})
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessCreateRepo, nil), nil
}
