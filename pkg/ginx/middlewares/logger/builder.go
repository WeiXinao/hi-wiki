package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody  *atomic.Bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc:   fn,
		allowReqBody: atomic.NewBool(false),
	}
}

func (b *MiddlewareBuilder) AllowReqBody(ok bool) *MiddlewareBuilder {
	b.allowReqBody.Store(ok)
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024] + "..."
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			//	URL 本身也可能很长
			Url: url,
		}
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			//	Body 读完就没有了，所以读出来后，要放回去，让 后面的 handler 可以拿到
			body, _ := ctx.GetRawData()
			reqBody := make([]byte, len(body))
			copy(reqBody, body)
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader
			if len(reqBody) > 1024 {
				reqBody = append(reqBody[:1024], []byte("...")...)
			}

			// 这其实是一个很消耗 CPU 和内存的操作
			// 因为会引起复制
			al.ReqBody = string(reqBody)
		}

		if b.allowRespBody {
			ctx.Writer = responseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			b.loggerFunc(ctx, al)
		}()

		// 执行业务逻辑
		ctx.Next()
	}
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// 整个请求的 URL
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}
