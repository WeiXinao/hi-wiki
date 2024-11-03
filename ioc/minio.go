package ioc

import (
	"github.com/WeiXinao/hi-wiki/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio(appConfig *config.AppConfig) *minio.Client {
	minioClient, err := minio.New(appConfig.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(appConfig.Minio.AccessKeyID, appConfig.Minio.SecureAccessKey, ""),
		Secure: appConfig.Minio.UseSSL,
	})
	if err != nil {
		panic(err)
	}
	return minioClient
}
