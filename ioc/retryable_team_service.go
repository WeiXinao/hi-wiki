package ioc

import (
	"github.com/WeiXinao/hi-wiki/internal/repository"
	"github.com/WeiXinao/hi-wiki/internal/service"
)

func InitRetryableTeamService(repo repository.TeamRepository) service.TeamService {
	retryMax := 3
	return service.NewRetryableTeamService(
		service.NewTeamService(repo),
		retryMax,
	)
}
