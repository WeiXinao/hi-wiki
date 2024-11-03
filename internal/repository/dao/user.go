package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrRecordNotFound        = gorm.ErrRecordNotFound
	ErrDuplicateUsername     = errors.New("用户名冲突")
	ErrPasswordHasBeenModify = errors.New("密码已经被修改了")
	ErrFailModifyProfile     = errors.New("修改个人信息失败")
)

type UserDao interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	Insert(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
	UpdateByIdAndPassword(ctx context.Context, id int64, oldPwd string, newPwd string) error
	UpdateById(ctx context.Context, avatarUrl string, desc string, uid int64) error
}

type gormUserDao struct {
	db *gorm.DB
}

func (g *gormUserDao) UpdateById(ctx context.Context, avatarUrl string, desc string, uid int64) error {
	res := g.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", uid).
		Updates(map[string]any{
			"avatar_md5": avatarUrl,
			"profile":    desc,
			"updated_at": time.Now().UnixMilli(),
		})
	resErr := res.Error
	if resErr != nil {
		return resErr
	}
	if res.RowsAffected == 0 {
		return ErrPasswordHasBeenModify
	}
	return nil
}

func (g *gormUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := g.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Find(&user).Error
	return user, err
}

func (g *gormUserDao) UpdateByIdAndPassword(ctx context.Context, id int64, oldPwd string, newPwd string) error {
	res := g.db.WithContext(ctx).Model(&User{}).Where("id = ? and password = ?", id, oldPwd).
		Updates(map[string]any{
			"password":   newPwd,
			"updated_at": time.Now().UnixMilli(),
		})
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return ErrPasswordHasBeenModify
	}
	return nil
}

func (g *gormUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreatedAt = now
	u.UpdatedAt = now
	u.DeletedAt = 0
	err := g.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrDuplicateUsername
		}
	}
	return err
}

func (g *gormUserDao) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := g.db.WithContext(ctx).
		Model(&User{}).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func NewUserDao(db *gorm.DB) UserDao {
	return &gormUserDao{db: db}
}

type User struct {
	BaseModel
	Username  sql.NullString `gorm:"unique"`
	Password  string
	AvatarMd5 string
	Profile   string
}
