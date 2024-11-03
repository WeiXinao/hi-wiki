package dao

import (
	"context"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type FileDao interface {
	InsertFile(ctx context.Context, file File) (File, error)
	InsertOwner(ctx context.Context, uid int64, md5Str string) (File, error)
	GetImageByMd5(ctx context.Context, md5Str string) (File, error)
}

var ErrImageNotFound = gorm.ErrRecordNotFound

type fileDao struct {
	db *gorm.DB
}

func (f *fileDao) GetImageByMd5(ctx context.Context, md5Str string) (File, error) {
	var file File
	err := f.db.WithContext(ctx).Model(&File{}).Where("md5 = ?", md5Str).First(&file).Error
	return file, err
}

func (f *fileDao) InsertOwner(ctx context.Context, uid int64, md5 string) (File, error) {
	var (
		file File
	)
	err := f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 通过 md5 找到 file
		err := tx.WithContext(ctx).Model(&File{}).Where("md5 = ?", md5).First(&file).Error
		if err != nil {
			return err
		}

		now := time.Now().UnixMilli()
		owner := Owner{
			BaseModel: BaseModel{
				Id:        0,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: 0,
			},
			UserId: uid,
			FileId: file.Id,
		}
		// 插入记录
		err = tx.WithContext(ctx).Create(&owner).Error
		if err == nil {
			return nil
		}
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			const uniqueConflictsErrNo uint16 = 1062
			if mysqlErr.Number == uniqueConflictsErrNo {
				return nil
			}
		}
		return err
	})
	return file, err
}

func (f *fileDao) InsertFile(ctx context.Context, file File) (File, error) {
	now := time.Now().UnixMilli()
	file.CreatedAt = now
	file.UpdatedAt = now
	file.DeletedAt = 0
	err := f.db.WithContext(ctx).Create(&file).Error
	return file, err
}

func NewFileDao(db *gorm.DB) FileDao {
	return &fileDao{
		db: db,
	}
}

type File struct {
	BaseModel
	Name           string
	Typ            string
	Md5            string `gorm:"unique"`
	Url            string
	Size           string
	DirLevel       uint
	RepoUniqueCode string `gorm:"type:varchar(191)"`
	UserId         int64
}

// Owner 什么用户拥有什么文件
type Owner struct {
	BaseModel
	// 在 UserId 和 FileId 上创建联合唯一索引
	UserId int64 `gorm:"unique_index:uk_uid_fid"`
	FileId int64 `gorm:"unique_index:uk_uid_fid"`
}
