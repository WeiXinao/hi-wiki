package oss

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

type OSS interface {
	PutObject(ctx context.Context, objectName string, reader io.Reader,
		objectSize int64, options minio.PutObjectOptions) (minio.UploadInfo, error)
	PutFileObject(ctx context.Context, objectName string, filePath string,
		opts minio.PutObjectOptions) (minio.UploadInfo, error)
	GetObject(ctx context.Context, objectName string, opts minio.GetObjectOptions) (*minio.Object, error)
	RemoveObject(ctx context.Context, objectName string,
		opts minio.RemoveObjectOptions) error
}

type minioOSS struct {
	minioClient *minio.Client
	bucketName  string
}

func (m *minioOSS) RemoveObject(ctx context.Context, objectName string,
	opts minio.RemoveObjectOptions) error {
	return m.minioClient.RemoveObject(ctx, m.bucketName, objectName, opts)
}

func (m *minioOSS) PutObject(ctx context.Context, objectName string, reader io.Reader,
	objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return m.minioClient.PutObject(ctx, m.bucketName, objectName, reader, objectSize, opts)
}

func (m *minioOSS) PutFileObject(ctx context.Context, objectName string, filePath string,
	opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	return m.minioClient.FPutObject(ctx, m.bucketName, objectName, filePath, opts)
}

func (m *minioOSS) GetObject(ctx context.Context, objectName string,
	opts minio.GetObjectOptions) (*minio.Object, error) {
	return m.minioClient.GetObject(ctx, m.bucketName, objectName, opts)
}

func NewMinioOSS(minioClient *minio.Client) OSS {
	return &minioOSS{
		minioClient: minioClient,
		bucketName:  "hi-wiki",
	}
}
