package dao

import (
	"context"
	"errors"
	"github.com/WeiXinao/xkit/slice"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicatedRepoUniqueCode = errors.New("知识库特征码重复")
)

type RepoDao interface {
	InsertRepo(ctx context.Context, repo Repo) error
	GetByUserId(ctx context.Context, uid int64, isDoc bool) ([]Repo, error)
	GetByGroupId(ctx context.Context, groupId int64) ([]Repo, error)
	GetByUniqueCode(ctx context.Context, code string) (Repo, error)
	GetHotDocRepo(ctx context.Context, doc bool, hot bool) ([]Repo, error)
}

type gormRepoDao struct {
	db *gorm.DB
}

func (g *gormRepoDao) GetHotDocRepo(ctx context.Context, isDoc bool, isHot bool) ([]Repo, error) {
	type UcAndCnt struct {
		RepoUniqueCode string
		Count          int64
	}
	var (
		ucAndCnts []UcAndCnt
		repoList  []Repo
	)
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Article{}).
			Select("repo_unique_code, count(like_cnt) AS c").
			Group("repo_unique_code").
			Order("c DESC").
			Limit(10).
			Find(&ucAndCnts).
			Error
		if err != nil {
			return err
		}
		codes := slice.Map[UcAndCnt, string](ucAndCnts, func(idx int, src UcAndCnt) string {
			return src.RepoUniqueCode
		})
		return tx.Model(&Repo{}).
			Where("unique_code IN ? AND status = ?", codes, 1).
			Find(&repoList).Error
	})
	return repoList, err
}

func (g *gormRepoDao) GetByUniqueCode(ctx context.Context, code string) (Repo, error) {
	var repo Repo
	err := g.db.WithContext(ctx).Model(&Repo{}).Where("unique_code = ?", code).First(&repo).Error
	return repo, err
}

func (g *gormRepoDao) GetByGroupId(ctx context.Context, groupId int64) ([]Repo, error) {
	var repoList []Repo
	err := g.db.WithContext(ctx).Where("team_id = ?", groupId).
		Preload("RepoCate").
		Find(&repoList).Error
	if err != nil {
		return nil, err
	}
	return repoList, nil
}

func (g *gormRepoDao) GetByUserId(ctx context.Context, uid int64, isDoc bool) ([]Repo, error) {
	var (
		repoList     []Repo
		selfRepoList []Repo
		err          error
	)
	if isDoc {
		err = g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// 获取团队知识库
			err := tx.InnerJoins("join teams on teams.id = repos.team_id").
				InnerJoins("join team_members on team_members.team_unique_code = teams.unique_code").
				Where("team_members.user_id = ? AND repos.cate = ?", uid, 1).
				Preload("RepoCate").
				Find(&repoList).Error
			if err != nil {
				return err
			}
			//	获取个人知识库
			return tx.Where("user_id = ? AND cate = ? AND team_id = ? ", uid, 1, 9999).
				Preload("RepoCate").
				Find(&selfRepoList).Error
		})
	} else {
		err = g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// 获取团队知识库
			err := tx.InnerJoins("join teams on teams.id = repos.team_id").
				InnerJoins("join team_members on team_members.team_unique_code = teams.unique_code").
				Where("team_members.user_id = ?", uid).
				Preload("RepoCate").
				Find(&repoList).Error
			if err != nil {
				return err
			}
			//	获取个人知识库
			return tx.Where("user_id = ? and team_id = ?", uid, 9999).
				Preload("RepoCate").
				Find(&selfRepoList).Error
		})
	}
	if err != nil {
		return nil, err
	}
	repoList = append(repoList, selfRepoList...)
	return repoList, nil
}

func (g *gormRepoDao) InsertRepo(ctx context.Context, repo Repo) error {
	now := time.Now().UnixMilli()
	repo.CreatedAt = now
	repo.UpdatedAt = now
	repo.DeletedAt = 0
	err := g.db.WithContext(ctx).Create(&repo).Error
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			const uniqueConflictsErrNo uint16 = 1062
			if mysqlErr.Number == uniqueConflictsErrNo {
				return ErrDuplicatedTeamUniqueCode
			}
		}
	}
	return err
}

func NewGormRepoDao(db *gorm.DB) RepoDao {
	return &gormRepoDao{
		db: db,
	}
}

type Repo struct {
	BaseModel
	Name       string
	Desc       string
	Cate       int64
	Status     uint8
	UniqueCode string `gorm:"unique;type:varchar(191)"`
	State      uint8
	TeamId     int64
	UserId     int64
	RepoCate   RepoCate `gorm:"foreignKey:Cate"`
}

type RepoCate struct {
	BaseModel
	Name string
}
