package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"github.com/WeiXinao/xkit/slice"
	"io"
	"mime/multipart"
	"os"
	path2 "path"
	"strings"
	"time"
)

var (
	ErrFailUploadFile = errors.New("上传文件失败")
	ErrIllegalPath    = errors.New("非法路径")
	ErrImageNotFound  = repository.ErrImageNotFound
)

type FileService interface {
	Upload(ctx context.Context, file *multipart.FileHeader, code string, level uint, uid int64) (domain.File, error)
	ShowImage(ctx context.Context, md5Str string) (domain.File, error)
	IsImage(ext string) bool
}

type fileService struct {
	repo     repository.FileRepository
	imgTypes []string
}

func (f *fileService) ShowImage(ctx context.Context, md5Str string) (domain.File, error) {
	file, err := f.repo.GetImageByMd5(ctx, md5Str)
	if err != nil {
		return domain.File{}, err
	}
	exists, err := f.fileExists(file.Url)
	if err != nil {
		return domain.File{}, err
	}
	if !exists {
		return domain.File{}, ErrImageNotFound
	}
	return file, nil
}

func (f *fileService) IsImage(ext string) bool {
	return slice.Contains[string](f.imgTypes, ext)
}

func (f *fileService) Upload(ctx context.Context, file *multipart.FileHeader, code string, level uint, uid int64) (domain.File, error) {
	var (
		filepath = ""
	)
	filenames := strings.Split(file.Filename, ".")
	if len(filenames) <= 1 || !f.IsImage(filenames[1]) {
		filepath += "file/"
	} else {
		filepath += "img/"
	}

	size := f.formatFileSize(file.Size)
	ext := ""
	if len(filenames) == 2 {
		ext = filenames[1]
	}
	hasher := md5.New()
	src, err := file.Open()
	if err != nil {
		return domain.File{}, err
	}
	// 计算文件的 md5 值
	_, err = io.Copy(hasher, src)
	if err != nil {
		return domain.File{}, err
	}
	md5Str := hex.EncodeToString(hasher.Sum(nil))

	// 拼接文件的完整路径
	nowDate := time.Now().Format("2006-01-02")
	parentPath := filepath + nowDate
	err = f.makeDir(parentPath)
	if err != nil {
		return domain.File{}, err
	}
	fullPath := strings.Join([]string{parentPath, md5Str}, "/")
	fullPath = fmt.Sprintf("%s.%s", fullPath, ext)
	// 有并发问题，暂且不处理
	exists, err := f.fileExists(fullPath)
	if err != nil {
		return domain.File{}, err
	}
	if !exists {
		//	不存在就上传
		if err := f.saveUploadedFile(file, fullPath); err != nil {
			return domain.File{}, ErrFailUploadFile
		}
		fileDomain := domain.File{
			Name:           file.Filename,
			Typ:            ext,
			Md5:            md5Str,
			Url:            fullPath,
			Size:           size,
			DirLevel:       level,
			RepoUniqueCode: code,
			Owner: domain.Owner{
				Id: uid,
			},
		}
		fileDomain, err = f.repo.InsertFile(ctx, fileDomain)
		//	将文件信息添加到数据库失败了怎么办？
		//	打日志，Prometheus 打点，重试 3 次，还不对，手工介入
		if err != nil {
			return domain.File{}, fmt.Errorf("添加文件信息到数据库失败: %w", err)
		}
	}

	//	文件已经存在，将它添加到当前用户名下
	return f.repo.InsertOwner(ctx, uid, md5Str)
}

func (f *fileService) saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// 创建目录
func (f *fileService) makeDir(path string) error {
	exists, err := f.exists(path)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	err = os.MkdirAll(path, 0755)
	return err
}

func (f *fileService) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断文件是否存在
func (f *fileService) fileExists(path string) (bool, error) {
	pathSeg := strings.Split(path, "/")
	if len(pathSeg) != 3 {
		return false, ErrIllegalPath
	}
	dirs, err := os.ReadDir(pathSeg[0])
	if err != nil {
		return false, err
	}
	var b bool
	for _, dir := range dirs {
		b, err = f.exists(path2.Join(pathSeg[0], dir.Name(), pathSeg[2]))
		if b {
			return true, nil
		}
	}
	return false, err
}

func (f *fileService) formatFileSize(size int64) string {
	const (
		B = 1 << (10 * iota)
		KB
		MB
		GB
		TB
		EB
	)
	switch {
	case size < KB:
		return fmt.Sprintf("%.2fB", float64(size)/float64(B))
	case size < MB:
		return fmt.Sprintf("%.2fKB", float64(size)/float64(KB))
	case size < GB:
		return fmt.Sprintf("%.2fMB", float64(size)/float64(MB))
	case size < TB:
		return fmt.Sprintf("%.2fGB", float64(size)/float64(GB))
	case size < EB:
		return fmt.Sprintf("%.2fTB", float64(size)/float64(TB))
	default:
		return fmt.Sprintf("%.2fEB", float64(size)/float64(EB))
	}
}

func NewFileService(repo repository.FileRepository) FileService {
	var imgType = []string{"jpg", "png", "gif", "webp", "svg", "apng", "jpeg"}
	return &fileService{
		repo:     repo,
		imgTypes: imgType,
	}
}
