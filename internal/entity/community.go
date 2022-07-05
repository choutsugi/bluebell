package entity

import (
	"gorm.io/gorm"
	"time"
)

type Community struct {
	Cid       uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(32);unique;not null" json:"name"`
	Intro     string `gorm:"type:varchar(32);unique;not null" json:"intro"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
