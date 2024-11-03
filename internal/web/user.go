package web

import (
	"errors"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/internal/service"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/internal/web/vo"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc service.UserService
	_jwt.HandlerJWT
}

func NewUserHandler(svc service.UserService, jwtHdl _jwt.HandlerJWT) *UserHandler {
	return &UserHandler{
		svc:        svc,
		HandlerJWT: jwtHdl,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/login", ginx.WrapBody[LoginReq](h.Login))
	ug.POST("/signup", ginx.WrapBody[SignUpReq](h.SignUp))
	ug.POST("/mod_pwd", ginx.WrapBodyAndClaims[ModifyPasswordReq, *_jwt.UserClaims](h.ModifyPassword))
	ug.POST("/mod_profile", ginx.WrapBodyAndClaims[ProfileReq, *_jwt.UserClaims](h.ModifyProfile))
	ug.GET("/profile", ginx.WrapClaims[*_jwt.UserClaims](h.Profile))
}

func (h *UserHandler) Profile(ctx *gin.Context, claims *_jwt.UserClaims) (Result, error) {
	user, err := h.svc.Profile(ctx.Request.Context(), claims.UserId)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessGetProfile, vo.UserVO{
		Id:        user.Id,
		Username:  user.Username,
		AvatarMd5: user.AvatarMd5,
		AvatarUrl: "files/img/" + user.AvatarMd5,
		Profile:   user.Profile,
	}), nil
}

type ProfileReq struct {
	AvatarUrl string `json:"avatarurl"`
	UserDesc  string `json:"userdesc"`
}

func (h *UserHandler) ModifyProfile(ctx *gin.Context, req ProfileReq, userClaims *_jwt.UserClaims) (Result, error) {
	err := h.svc.ModifyProfile(ctx.Request.Context(), req.AvatarUrl, req.UserDesc, userClaims.UserId)
	if errors.Is(err, service.ErrFailModifyProfile) {
		return MarshalResp(errs.FailModifyProfile, nil), nil
	}
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}

	return MarshalResp(errs.SuccessModifyProfile, nil), nil
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(ctx *gin.Context, req LoginReq) (Result, error) {
	u, err := h.svc.Login(ctx.Request.Context(), req.Username, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		return MarshalResp(errs.InvalidUserOrPassword, nil), nil
	}
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), nil
	}

	if err = h.SetJwtToken(ctx, u.Id); err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}

	return MarshalResp(errs.SuccessLogin, nil), nil
}

type SignUpReq struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	AgainPassword string `json:"againPassword"`
}

func (h *UserHandler) SignUp(ctx *gin.Context, req SignUpReq) (Result, error) {
	// 参数校验
	if len(req.Username) < 5 || len(req.Username) > 20 {
		return MarshalResp(errs.UsernameOutOfRange, nil), nil
	}
	if req.Password != req.AgainPassword {
		return MarshalResp(errs.InconsistentPasswordAndConfirmPassword, nil), nil
	}
	if len(req.Password) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}

	err := h.svc.SignUp(ctx.Request.Context(), domain.User{
		Username: req.Username,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrDuplicateUsername) {
		return MarshalResp(errs.UserDuplicatedUsername, nil), nil
	}
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}

	return MarshalResp(errs.SuccessSign, nil), nil
}

type ModifyPasswordReq struct {
	OldPasswd    string `json:"old_passwd,omitempty"`
	NewPasswdOne string `json:"new_passwd_one,omitempty"`
	NewPasswdTwo string `json:"new_passwd_two,omitempty"`
}

func (h *UserHandler) ModifyPassword(ctx *gin.Context, req ModifyPasswordReq, userClaims *_jwt.UserClaims) (Result, error) {
	//	参数校验
	if req.NewPasswdOne != req.NewPasswdTwo {
		return MarshalResp(errs.InconsistentTwoPassword, nil), nil
	}

	err := h.svc.ModifyPassword(ctx.Request.Context(), userClaims.UserId, req.OldPasswd, req.NewPasswdOne)
	if errors.Is(err, service.ErrInvalidPassword) || errors.Is(err, service.ErrPasswordHasBeenModify) {
		return MarshalResp(errs.InvalidPassword, nil), nil
	}

	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessModifyPassword, nil), nil
}
