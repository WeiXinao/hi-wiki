package dao

type BaseModel struct {
	Id        int64 `gorm:"primaryKey,autoIncrement"`
	CreatedAt int64
	UpdatedAt int64
	DeletedAt int64 `gorm:"default:0"`
}
