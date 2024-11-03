package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&File{},
		&Owner{},
		&Team{},
		&TeamMember{},
		&Repo{},
		&RepoCate{},
		&Article{},
		&GoodHistory{},
		&Follow{},
		&BookCate{},
		&Book{},
	)
}
