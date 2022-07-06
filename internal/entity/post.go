package entity

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID          int64  `gorm:"primarykey"`
	Title       string `gorm:"type:varchar(128);not null" json:"title"`
	Content     string `gorm:"type:text;not null" json:"content"`
	AuthorId    int64  `gorm:"type:int;not null" json:"author_id"`
	CommunityId int64  `gorm:"type:int;not null" json:"community_id"`
	Status      int    `gorm:"type:int;not null" json:"status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
