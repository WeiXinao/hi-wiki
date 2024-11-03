package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository"
)

var (
	ErrDuplicatedTeamUniqueCode  = repository.ErrDuplicatedTeamUniqueCode
	ErrTeamLeaderCannotBeDeleted = repository.ErrTeamLeaderCannotBeDeleted
)

type TeamService interface {
	Create(ctx context.Context, name string, desc string, auth uint8, avatar string, uid int64) (string, error)
	ListTeams(ctx context.Context, uid int64) ([]domain.Team, error)
	Detail(ctx context.Context, code string) (domain.Team, error)
	ListTeamMembers(ctx context.Context, code string, uid int64) ([]domain.TeamMember, error)
	DeleteTeamMember(ctx context.Context, code string, uid int64) error
}

type teamService struct {
	repo repository.TeamRepository
}

func (t *teamService) DeleteTeamMember(ctx context.Context, code string, uid int64) error {
	return t.repo.DeleteTeamMember(ctx, code, uid)
}

func (t *teamService) ListTeamMembers(ctx context.Context, code string, uid int64) ([]domain.TeamMember, error) {
	return t.repo.GetTeamMembers(ctx, code, uid)
}

func (t *teamService) Detail(ctx context.Context, code string) (domain.Team, error) {
	return t.repo.GetTeamByUniqueCode(ctx, code)
}

func (t *teamService) ListTeams(ctx context.Context, uid int64) ([]domain.Team, error) {
	return t.repo.GetTeams(ctx, uid)
}

func (t *teamService) Create(ctx context.Context, name string, desc string, auth uint8, avatar string, uid int64) (string, error) {
	uniqueCode := t.generateRandomString(32)
	// 传入的是头像的 md5 值
	err := t.repo.InsertTeamAndMember(ctx, name, desc, uniqueCode, auth, avatar, uid)
	if err != nil {
		return "", err
	}
	return uniqueCode, nil
}

func (t *teamService) generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func NewTeamService(repo repository.TeamRepository) TeamService {
	return &teamService{
		repo: repo,
	}
}
