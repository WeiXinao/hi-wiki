package repository

import (
	"context"
	"github.com/WeiXinao/hi-wiki/internal/domain"
	"github.com/WeiXinao/hi-wiki/internal/repository/cache"
	"github.com/WeiXinao/hi-wiki/internal/repository/dao"
	"github.com/WeiXinao/xkit/slice"
	"time"
)

var (
	ErrDuplicatedTeamUniqueCode  = dao.ErrDuplicatedTeamUniqueCode
	ErrTeamLeaderCannotBeDeleted = dao.ErrTeamLeaderCannotBeDeleted
)

type TeamRepository interface {
	InsertTeamAndMember(ctx context.Context, name string, desc string, code string, auth uint8, avatar string, uid int64) error
	GetTeams(ctx context.Context, uid int64) ([]domain.Team, error)
	GetTeamByUniqueCode(ctx context.Context, code string) (domain.Team, error)
	GetTeamMembers(ctx context.Context, code string, uid int64) ([]domain.TeamMember, error)
	DeleteTeamMember(ctx context.Context, code string, uid int64) error
}

type teamRepository struct {
	dao   dao.TeamDao
	cache cache.TeamCache
}

func (t *teamRepository) DeleteTeamMember(ctx context.Context, code string, uid int64) error {
	err := t.dao.DeleteTeamMember(ctx, code, uid)

	go func() {
		_ = t.cache.DelTeamMember(ctx, code, uid)
	}()

	return err
}

func (t *teamRepository) GetTeamMembers(ctx context.Context, code string, uid int64) ([]domain.TeamMember, error) {
	members, err := t.cache.GetTeamMembers(ctx, code)
	if err == nil {
		return members, nil
	}

	memberModels, err := t.dao.GetTeamMembers(ctx, code, uid)
	if err != nil {
		return nil, err
	}
	members = slice.Map[dao.TeamMember, domain.TeamMember](memberModels, func(idx int, src dao.TeamMember) domain.TeamMember {
		return domain.TeamMember{
			Id:        src.Id,
			Name:      src.User.Username.String,
			IsLeader:  src.IsLeader,
			UserId:    src.UserId,
			CreateAt:  time.UnixMilli(src.CreatedAt),
			UpdatedAt: time.UnixMilli(src.UpdatedAt),
			DeletedAt: time.UnixMilli(src.DeletedAt),
		}
	})
	go func() {
		_ = t.cache.AddTeamMembers(ctx, code, members...)
	}()
	return members, err
}

func (t *teamRepository) GetTeamByUniqueCode(ctx context.Context, code string) (domain.Team, error) {
	team, err := t.cache.GetTeam(ctx, code)
	if err == nil {
		return team, nil
	}

	teamModel, err := t.dao.GetTeamByUniqueCode(ctx, code)
	team = t.entityToDomain(teamModel)

	go func() {
		_ = t.cache.SetTeam(ctx, team)
	}()

	return team, err
}

func (t *teamRepository) GetTeams(ctx context.Context, uid int64) ([]domain.Team, error) {
	teams, err := t.dao.GetTeams(ctx, uid)
	return slice.Map[dao.Team, domain.Team](teams, func(idx int, src dao.Team) domain.Team {
		return t.entityToDomain(src)
	}), err
}

func (t *teamRepository) InsertTeamAndMember(ctx context.Context, name string, desc string,
	code string, auth uint8, avatar string, uid int64) error {
	return t.dao.InsertTeamAndMember(ctx, name, desc, code, auth, avatar, uid)
}

func NewTeamRepository(dao dao.TeamDao, cache cache.TeamCache) TeamRepository {
	return &teamRepository{
		dao:   dao,
		cache: cache,
	}
}

func (t *teamRepository) entityToDomain(team dao.Team) domain.Team {
	return domain.Team{
		Id:         team.Id,
		Name:       team.Name,
		Desc:       team.Desc,
		UniqueCode: team.UniqueCode,
		Status:     domain.TeamStatus(team.Status),
		AvatarMd5:  team.AvatarMd5,
		CreateAt:   time.UnixMilli(team.CreatedAt),
		UpdatedAt:  time.UnixMilli(team.UpdatedAt),
		DeletedAt:  time.UnixMilli(team.DeletedAt),
	}
}

func (t *teamRepository) domainToEntity(team domain.Team) dao.Team {
	return dao.Team{
		BaseModel: dao.BaseModel{
			Id:        team.Id,
			CreatedAt: team.CreateAt.UnixMilli(),
			UpdatedAt: team.UpdatedAt.UnixMilli(),
			DeletedAt: team.DeletedAt.UnixMilli(),
		},
		Name:       team.Name,
		Desc:       team.Desc,
		UniqueCode: team.UniqueCode,
		Status:     team.Status.ToUint8(),
		AvatarMd5:  team.AvatarMd5,
	}
}
