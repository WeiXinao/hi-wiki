package web

import (
	"errors"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/internal/service"
	_jwt "github.com/WeiXinao/hi-wiki/internal/web/jwt"
	"github.com/WeiXinao/hi-wiki/internal/web/vo"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
	"github.com/WeiXinao/xkit/slice"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type BookHandler struct {
	svc     service.BookService
	fileSvc service.FileService
}

func NewBookHandler(svc service.BookService, fileSvc service.FileService) *BookHandler {
	return &BookHandler{
		svc:     svc,
		fileSvc: fileSvc,
	}
}

func (h *BookHandler) RegisterRoutes(server *gin.Engine) {
	bg := server.Group("/books")
	bg.POST("/cate/edit", ginx.WrapBody[BookCateReq](h.EditCate))
	bg.GET("/cate/list", ginx.Wrap(h.ListCate))
	bg.DELETE("/cate/del/:cid", ginx.Wrap(h.DelCate))
	bg.POST("/upload", ginx.WrapClaims[*_jwt.UserClaims](h.UploadBook))
	bg.GET("/list", ginx.Wrap(h.List))
	bg.GET("/download/:id", h.Download)
}

func (h *BookHandler) Download(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, MarshalResp(errs.InternalInvalidInput, nil))
		return
	}
	book, err := h.svc.Download(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusOK, MarshalResp(errs.InternalServerError, nil))
		return
	}
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", book.BookName))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(book.BookUrl)
}

func (h *BookHandler) List(ctx *gin.Context) (ginx.Result, error) {
	cate := ctx.DefaultQuery("flag", "all")
	if strings.TrimSpace(cate) == "all" {
		cate = "0"
	}
	// 关键词
	kw := ctx.Query("kw")
	// 排序
	rank := ctx.Query("rank")
	cateInt64, err := strconv.ParseInt(cate, 10, 64)
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	books, err := h.svc.List(ctx.Request.Context(), cateInt64, kw, rank, page, size)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessListBook, slice.Map[domain.Book, vo.BookVO](books,
		func(idx int, src domain.Book) vo.BookVO {
			return vo.BookVO{
				Id:            src.Id,
				BookName:      src.BookName,
				BookUrl:       src.BookUrl,
				BookAvatarUrl: src.GetAvatarUrl(),
				BookCateName:  src.BookCateName,
				Download:      src.Download,
			}
		})), nil
}

func (h *BookHandler) UploadBook(ctx *gin.Context, userClaims *_jwt.UserClaims) (Result, error) {
	type fileInfo struct {
		file           *multipart.FileHeader
		repoUniqueCode string
	}
	var (
		dirLevel  uint = 0
		fileInfos []fileInfo
		fileUrls  []string
		md5s      []string
	)
	_ = ctx.PostForm("filetype")
	filecate := ctx.PostForm("filecate")
	file, err := ctx.FormFile("file")
	avatar, err := ctx.FormFile("avator")
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), err
	}
	fileInfos = append(fileInfos, fileInfo{
		file:           file,
		repoUniqueCode: "book",
	})
	fileInfos = append(fileInfos, fileInfo{
		file:           avatar,
		repoUniqueCode: "avatar",
	})

	//	参数校验
	filecateInt64, err := strconv.ParseInt(filecate, 10, 64)
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), err
	}
	for _, fi := range fileInfos {
		uploadedFile, err := h.fileSvc.Upload(ctx.Request.Context(), fi.file, fi.repoUniqueCode, dirLevel, userClaims.UserId)
		if errors.Is(err, service.ErrFailUploadFile) {
			return MarshalResp(errs.FailUploadFile, nil), nil
		}
		fileUrls = append(fileUrls, uploadedFile.Url)
		md5s = append(md5s, uploadedFile.Md5)
	}
	if len(fileUrls) == 2 && len(md5s) == 2 {
		err := h.svc.Add(ctx.Request.Context(), fileInfos[0].file.Filename, fileUrls[0], fileUrls[1], filecateInt64,
			userClaims.UserId, md5s[0], md5s[1])
		if err != nil {
			return MarshalResp(errs.FailAddBookInfo, nil), err
		}
	}

	return MarshalResp(errs.SuccessUploadBook, nil), nil
}

type BookCateReq struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (h *BookHandler) DelCate(ctx *gin.Context) (ginx.Result, error) {
	//	获取参数并校验
	id, err := strconv.ParseInt(ctx.Param("cid"), 10, 64)
	if err != nil {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	err = h.svc.DelCate(ctx.Request.Context(), id)
	if errors.Is(err, service.ErrDeletedBookCateNotFound) {
		return MarshalResp(errs.DeletedBookCateNotFound, nil), nil
	}
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessDeleteCate, nil), nil
}

func (h *BookHandler) ListCate(ctx *gin.Context) (ginx.Result, error) {
	bookCates, err := h.svc.ListCate(ctx.Request.Context())
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessListBookCates, slice.Map[domain.BookCate, vo.BookCateVO](bookCates, func(idx int, src domain.BookCate) vo.BookCateVO {
		return vo.BookCateVO{
			Id:       src.Id,
			CateName: src.Name,
		}
	})), nil
}

func (h *BookHandler) EditCate(ctx *gin.Context, req BookCateReq) (ginx.Result, error) {
	// 参数校验
	if len(strings.TrimSpace(req.Name)) == 0 {
		return MarshalResp(errs.InternalInvalidInput, nil), nil
	}
	//	执行逻辑
	err := h.svc.EditCate(ctx.Request.Context(), req.Id, req.Name)
	if err != nil {
		return MarshalResp(errs.InternalServerError, nil), err
	}
	return MarshalResp(errs.SuccessEditBookCate, nil), nil
}
