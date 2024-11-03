package web

import (
	"errors"
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
)

type TeamHandler struct {
	svc service.TeamService
}

func NewTeamHandler(svc service.TeamService) *TeamHandler {
	return &TeamHandler{
		svc: svc,
	}
}

func (h *TeamHandler) RegisterRoutes(server *gin.Engine) {
	tg := server.Group("/teams")
	tg.POST("/create", ginx.WrapBodyAndClaims[TeamReq, *_jwt.UserClaims](h.CreateTeam))
	tg.GET("/list", ginx.WrapClaims[*_jwt.UserClaims](h.ListTeams))
	tg.GET("/detail/:flag", ginx.Wrap(h.Detail))
	tg.GET("/members/:teamflag", ginx.WrapClaims[*_jwt.UserClaims](h.ListTeamMembers))
	tg.DELETE("/members/:teamflag/:uid", ginx.Wrap(h.DeleteTeamMember))
}

func (h *TeamHandler) DeleteTeamMember(ctx *gin.Context) (Result, error) {
	groupflag := ctx.Param("teamflag")
	uid := ctx.Param("uid")

	//	参数校验
	if len(strings.TrimSpace(groupflag)) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	uidInt64, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), err
	}

	err = h.svc.DeleteTeamMember(ctx.Request.Context(), groupflag, uidInt64)
	if errors.Is(err, service.ErrTeamLeaderCannotBeDeleted) {
		return MarshalResp(errs.TeamLeaderCannotBeDeleted, nil), nil
	}
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessDeleteTeamMember, nil), nil
}

func (h *TeamHandler) ListTeamMembers(ctx *gin.Context, uc *_jwt.UserClaims) (Result, error) {
	uniqueCode := ctx.Param("teamflag")
	// 参数校验
	if len(strings.TrimSpace(uniqueCode)) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil),
			errors.New("参数校验错误，无效的团队特征码")
	}
	members, err := h.svc.ListTeamMembers(ctx.Request.Context(), uniqueCode, uc.UserId)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessListTeamMember,
		slice.Map[domain.TeamMember, vo.TeamMemberVO](members, func(idx int, src domain.TeamMember) vo.TeamMemberVO {
			userType := "组员"
			if src.IsLeader {
				userType = "组长"
			}
			return vo.TeamMemberVO{
				Uid:      src.UserId,
				Username: src.Name,
				UserType: userType,
			}
		})), nil
}

func (h *TeamHandler) Detail(ctx *gin.Context) (Result, error) {
	uniqueCode := ctx.Param("flag")
	// 参数校验
	if len(strings.TrimSpace(uniqueCode)) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	detail, err := h.svc.Detail(ctx.Request.Context(), uniqueCode)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessShowDetail, vo.TeamInfoVO{
		Id:         detail.Id,
		Name:       detail.Name,
		UniqueCode: detail.UniqueCode,
		Desc:       detail.Desc,
		Status:     detail.Status.ToUint8(),
		AvatarMd5:  detail.AvatarMd5,
		CreatedAt:  detail.CreateAt,
	}), nil
}

// ListTeams 获取和用户相关的团队列表
func (h *TeamHandler) ListTeams(ctx *gin.Context, uc *_jwt.UserClaims) (Result, error) {
	teams, err := h.svc.ListTeams(ctx.Request.Context(), uc.UserId)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	teamVOs := slice.Map[domain.Team, vo.TeamVO](teams, func(idx int, src domain.Team) vo.TeamVO {
		return vo.TeamVO{
			Id:         src.Id,
			Name:       src.Name,
			UniqueCode: src.UniqueCode,
			Desc:       src.Desc,
			AvatarMd5:  src.AvatarMd5,
		}
	})
	return MarshalResp(errs.SuccessListTeam, teamVOs), nil
}

type TeamReq struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Auth   uint8  `json:"auth"`
	Avatar string `json:"avator"`
}

// CreateTeam 创建团队
func (h *TeamHandler) CreateTeam(ctx *gin.Context, req TeamReq, uc *_jwt.UserClaims) (Result, error) {
	// 参数校验
	if !domain.RepoStatus(req.Auth).Valid() {
		return MarshalResp(errs.InternalInvalidInput, nil),
			errors.New("参数校验错误，无效的可见范围")
	}

	uniqueCode, err := h.svc.Create(ctx.Request.Context(), req.Name, req.Desc, req.Auth, req.Avatar, uc.UserId)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessCreateTeam, uniqueCode), err
}
