package dao

import (
	"context"
	"errors"
	"github.com/WeiXinao/xkit/slice"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrFollowedObjectNotExists = errors.New("被关注的对象不存在")

type FollowDao interface {
	GetFollowInfo(ctx context.Context, uids []int64,
		currentUser int64, followedIds []int64, typ uint8) (map[int64]FollowInfo, error)
	Insert(ctx context.Context, typ uint8, followedId int64, userId int64) error
	GetFollowedArticleAndRepo(ctx context.Context, uid int64) ([]FollowRecent, error)
}

type gormFollowDao struct {
	db *gorm.DB
}

func (g *gormFollowDao) GetFollowedArticleAndRepo(ctx context.Context, uid int64) ([]FollowRecent, error) {
	var (
		follows     []Follow
		artByUsers  []Article
		artsByRepos []Article
		repoIds     = make([]int64, 0, 10)
		userIds     = make([]int64, 0, 10)
	)
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("user_id = ?", uid).
			Find(&follows).Error
		if err != nil {
			return err
		}
		for _, follow := range follows {
			if follow.FollowType == 0 {
				userIds = append(userIds, follow.FollowId)
			} else {
				repoIds = append(repoIds, follow.FollowId)
			}
		}
		threeDaysBefore := time.Now().Add(-3 * 24 * time.Hour).UnixMilli()
		err = tx.Where("user_id IN ? AND created_at >= ?", userIds, threeDaysBefore).
			Preload("User").
			Find(&artByUsers).Error
		if err != nil {
			return err
		}
		return tx.Model(&Article{}).
			InnerJoins("JOIN repos ON articles.repo_unique_code = repos.unique_code").
			Where("repos.id IN ? AND articles.created_at >= ?", repoIds, threeDaysBefore).
			Preload("User").
			Find(&artsByRepos).Error
	})
	if err != nil {
		return nil, err
	}
	resByUsers := slice.Map[Article, FollowRecent](artByUsers, func(idx int, src Article) FollowRecent {
		return FollowRecent{
			UserName:   src.User.Username.String,
			Avatar:     src.User.AvatarMd5,
			UniqueCode: src.UniqueCode,
			Title:      src.Title,
			CreateTime: src.CreatedAt,
			Typ:        0,
		}
	})
	resByRepos := slice.Map[Article, FollowRecent](artsByRepos, func(idx int, src Article) FollowRecent {
		return FollowRecent{
			UserName:   src.User.Username.String,
			Avatar:     src.User.AvatarMd5,
			UniqueCode: src.UniqueCode,
			Title:      src.Title,
			CreateTime: src.CreatedAt,
			Typ:        1,
		}
	})
	return append(resByUsers, resByRepos...), nil
}

func (g *gormFollowDao) Insert(ctx context.Context, typ uint8, followedId int64,
	userId int64) error {
	var (
		user User
		repo Repo
		err  error
		now  = time.Now().UnixMilli()
	)
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		clauses := tx.Clauses(clause.Locking{Strength: "UPDATE"})
		if typ == 0 {
			err = clauses.First(&user).Error
		} else {
			err = clauses.First(&repo).Error
		}
		if errors.Is(err, ErrRecordNotFound) {
			return ErrFollowedObjectNotExists
		}
		if err != nil {
			return err
		}

		return tx.Create(&Follow{
			BaseModel: BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: 0,
			},
			FollowType: typ,
			FollowId:   followedId,
			UserId:     userId,
		}).Error
	})
}

func (g *gormFollowDao) GetFollowInfo(ctx context.Context, uids []int64,
	currentUser int64, followedIds []int64, typ uint8) (map[int64]FollowInfo, error) {
	var (
		followCounts []struct {
			UserId int64
			Count  int64
		}
		followedCounts []struct {
			FollowId int64
			Count    int64
		}
		follows []Follow
	)
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := g.db.Model(&Follow{}).
			Select("user_id, count(*) AS count").
			Group("user_id").
			Having("user_id IN ?", uids).
			Find(&followCounts).Error
		if err != nil {
			return err
		}
		err = g.db.Model(&Follow{}).
			Select("follow_id, count(*) AS count").
			Where("follow_type = ?", typ).
			Group("follow_id").
			Having("follow_id IN ?", followedIds).
			Find(&followedCounts).Error
		if err != nil {
			return err
		}
		return g.db.Model(&Follow{}).
			Where("follow_type = ? AND follow_id IN ? and user_id = ?",
				typ, followedIds, currentUser).
			Find(&follows).Error

	})
	if err != nil {
		return nil, err
	}
	res := make(map[int64]FollowInfo)
	for _, followCount := range followCounts {
		res[followCount.UserId] = FollowInfo{
			FollowCount: followCount.Count,
		}
	}
	for _, followedCount := range followedCounts {
		followInfo, ok := res[followedCount.FollowId]
		if !ok {
			res[followedCount.FollowId] = FollowInfo{
				BeFollowCount: followedCount.Count,
			}
			continue
		}
		followInfo.BeFollowCount = followedCount.Count
		res[followedCount.FollowId] = followInfo
	}
	for _, id := range followedIds {
		isFollowed := slice.ContainsFunc(follows, func(src Follow) bool {
			return src.FollowType == typ && src.UserId == currentUser && src.FollowId == id
		})
		info, ok := res[id]
		if !ok {
			res[id] = FollowInfo{
				IsFollowed: isFollowed,
			}
			continue
		}
		info.IsFollowed = isFollowed
		res[id] = info
	}
	return res, nil
}

func NewGormFollowDao(db *gorm.DB) FollowDao {
	return &gormFollowDao{
		db: db,
	}
}

type FollowRecent struct {
	UserName   string
	Avatar     string
	UniqueCode string
	Title      string
	CreateTime int64
	Typ        uint8
}

type Follow struct {
	BaseModel
	FollowType uint8 `gorm:"unique_index:uk_type_follow_id_user_id"`
	FollowId   int64 `gorm:"unique_index:uk_type_follow_id_user_id"`
	UserId     int64 `gorm:"unique_index:uk_type_follow_id_user_id"`
}

type FollowInfo struct {
	FollowCount   int64
	BeFollowCount int64
	IsFollowed    bool
}
