package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Uid       int64  `gorm:"primarykey"`
	Username  string `gorm:"type:varchar(32);unique;not null" json:"username"`
	Email     string `gorm:"type:varchar(32);unique;not null" json:"email"`
	Password  []byte `gorm:"type:bytea;not null"`
	Salt      []byte `gorm:"type:bytea;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
