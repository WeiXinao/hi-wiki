package ioc

import (
	"github.com/WeiXinao/hi-wiki/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
