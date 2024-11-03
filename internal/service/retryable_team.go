package service

import (
	"context"
	"errors"
)

type retryableTeamService struct {
	TeamService
	// 最大重试次数
	retryMax int
}

func (t *retryableTeamService) Create(ctx context.Context, name string, desc string,
	auth uint8, avatar string, uid int64) (string, error) {
	uniqueCode, err := t.TeamService.Create(ctx, name, desc, auth, avatar, uid)
	for i := 0; i < t.retryMax && errors.Is(err, ErrDuplicatedTeamUniqueCode); i++ {
		uniqueCode, err = t.TeamService.Create(ctx, name, desc, auth, avatar, uid)
		if !errors.Is(err, ErrDuplicatedTeamUniqueCode) {
			return uniqueCode, err
		}
	}
	return uniqueCode, err
}

func NewRetryableTeamService(svc TeamService, retryMax int) TeamService {
	return &retryableTeamService{
		TeamService: svc,
		retryMax:    retryMax,
	}
}
