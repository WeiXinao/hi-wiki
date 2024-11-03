package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicatedTeamUniqueCode  = errors.New("团队特征码重复")
	ErrTeamLeaderCannotBeDeleted = errors.New("不能删除组长")
)

type TeamDao interface {
	InsertTeamAndMember(ctx context.Context, name string, desc string, code string, auth uint8, avatar string, uid int64) error
	GetTeams(ctx context.Context, uid int64) ([]Team, error)
	GetTeamByUniqueCode(ctx context.Context, code string) (Team, error)
	GetTeamMembers(ctx context.Context, code string, uid int64) ([]TeamMember, error)
	DeleteTeamMember(ctx context.Context, code string, uid int64) error
}

type teamDao struct {
	db *gorm.DB
}

func (t *teamDao) DeleteTeamMember(ctx context.Context, code string, uid int64) error {
	res := t.db.WithContext(ctx).
		Where("user_id = ? AND team_unique_code = ? AND is_leader = ?", uid, code, false).
		Delete(&TeamMember{})
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return ErrTeamLeaderCannotBeDeleted
	}
	return nil
}

func (t *teamDao) GetTeamMembers(ctx context.Context, code string, uid int64) ([]TeamMember, error) {
	var teamMembers []TeamMember
	err := t.db.WithContext(ctx).
		Where("team_unique_code = ? AND deleted_at = 0", code).
		Preload("User").
		Find(&teamMembers).Error
	return teamMembers, err
}

func (t *teamDao) GetTeamByUniqueCode(ctx context.Context, code string) (Team, error) {
	var team Team
	err := t.db.WithContext(ctx).
		Model(&Team{}).
		Where("unique_code = ?", code).
		First(&team).Error
	return team, err
}

func (t *teamDao) GetTeams(ctx context.Context, uid int64) ([]Team, error) {
	var teamList []Team
	err := t.db.WithContext(ctx).Model(&Team{}).
		Joins("left join team_members on teams.unique_code = team_members.team_unique_code").
		Where("team_members.user_id = ?", uid).
		Scan(&teamList).Error
	return teamList, err
}

func (t *teamDao) InsertTeamAndMember(ctx context.Context, name string, desc string,
	code string, auth uint8, avatar string, uid int64) error {
	now := time.Now().UnixMilli()
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&Team{
			BaseModel: BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: 0,
			},
			Name:       name,
			Desc:       desc,
			UniqueCode: code,
			Status:     auth,
			AvatarMd5:  avatar,
		}).Error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			const uniqueConflictsErrNo uint16 = 1062
			if mysqlErr.Number == uniqueConflictsErrNo {
				return ErrDuplicatedTeamUniqueCode
			}
		}
		return tx.Create(&TeamMember{
			BaseModel: BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
			UserId:         uid,
			IsLeader:       true,
			TeamUniqueCode: code,
		}).Error
	})
}

func NewTeamDao(db *gorm.DB) TeamDao {
	return &teamDao{
		db: db,
	}
}

type Team struct {
	BaseModel
	Name string
	Desc string
	// 创建唯一索引
	UniqueCode string `gorm:"unique"`
	Status     uint8
	AvatarMd5  string
}

type TeamMember struct {
	BaseModel
	UserId         int64  `gorm:"index:idx_team_unique_code_user_id"`
	TeamUniqueCode string `gorm:"index:idx_team_unique_code_user_id"`
	IsLeader       bool
	User           User `gorm:"foreignKey:UserId"`
}
