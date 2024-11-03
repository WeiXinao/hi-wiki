package web

import (
	"errors"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/internal/service"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type FileHandler struct {
	svc service.FileService
}

func NewFileHandler(svc service.FileService) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (h *FileHandler) RegisterRoutes(server *gin.Engine) {
	fg := server.Group("/files")
	fg.POST("/upload", ginx.WrapClaims[*_jwt.UserClaims](h.Upload))
	fg.GET("/img/:md5", h.ShowImage)
}

func (h *FileHandler) ShowImage(ctx *gin.Context) {
	md5Str := ctx.Param("md5")

	image, err := h.svc.ShowImage(ctx.Request.Context(), md5Str)
	if errors.Is(err, service.ErrImageNotFound) {
		ctx.JSON(http.StatusOK, MarshalResp(errs.ImageNotFound, nil))
	}
	ctx.File(image.Url)
}

func (h *FileHandler) Upload(ctx *gin.Context, userClaims *_jwt.UserClaims) (Result, error) {
	dirLevelStr := ctx.PostForm("file_dir_level")
	repoUniqueCode := ctx.PostForm("repo_unique_code")
	file, err := ctx.FormFile("file")
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), err
	}

	//	参数校验
	parseUint, err := strconv.ParseUint(dirLevelStr, 10, 32)
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), err
	}
	dirLevel := uint(parseUint)
	if len(strings.Trim(repoUniqueCode, " ")) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}

	uploadedFile, err := h.svc.Upload(ctx.Request.Context(), file, repoUniqueCode, dirLevel, userClaims.UserId)
	if errors.Is(err, service.ErrFailUploadFile) {
		return MarshalResp(errs.FailUploadFile, nil), nil
	}

	return MarshalResp(errs.SuccessUploadFile, h.jointImgPath(uploadedFile)), nil
}

func (h *FileHandler) jointImgPath(file domain.File) map[string]string {
	const biz = "files"
	md5Str := file.Md5
	respData := map[string]string{
		"md5": md5Str,
	}
	if h.svc.IsImage(file.Typ) {
		respData["url"] = fmt.Sprintf("/%s/img/%s", biz, md5Str)
	}
	return respData
}
