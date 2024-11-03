package service

import (
	"context"
	"errors"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserOrPassword = errors.New("账号或密码不对")
	ErrInvalidPassword       = errors.New("密码错误")
	ErrDuplicateUsername     = repository.ErrDuplicateUsername
	ErrPasswordHasBeenModify = repository.ErrPasswordHasBeenModify
	ErrFailModifyProfile     = repository.ErrFailModifyProfile
)

type UserService interface {
	Login(ctx context.Context, username, password string) (domain.User, error)
	SignUp(ctx context.Context, user domain.User) error
	ModifyPassword(ctx context.Context, id int64, oldPwd string, newPwd string) error
	ModifyProfile(ctx context.Context, avatarUrl string, desc string, uid int64) error
	Profile(ctx context.Context, uid int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func (u *userService) Profile(ctx context.Context, uid int64) (domain.User, error) {
	user, err := u.repo.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (u *userService) ModifyProfile(ctx context.Context, avatarUrl string, desc string, uid int64) error {
	return u.repo.UpdateById(ctx, avatarUrl, desc, uid)
}

// ModifyPassword 修改密码
func (u *userService) ModifyPassword(ctx context.Context, id int64, oldPwd string, newPwd string) error {
	user, err := u.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd))
	if err != nil {
		return ErrInvalidPassword
	}
	newPwdBytes, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return u.repo.UpdateByIdAndPassword(ctx, id, user.Password, string(newPwdBytes))
}

func (u *userService) SignUp(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return u.repo.Create(ctx, user)
}

func (u *userService) Login(ctx context.Context, username, password string) (domain.User, error) {
	user, err := u.repo.FindByUsername(ctx, username)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
