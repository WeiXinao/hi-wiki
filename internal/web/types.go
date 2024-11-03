package web

import (
	"github.com/WeiXinao/hi-wiki/internal/errs"
	"github.com/WeiXinao/hi-wiki/pkg/ginx"
)

type Result = ginx.Result

func MarshalResp(err errs.Err, data any) Result {
	return ginx.MarshalResp(err, data)
}
