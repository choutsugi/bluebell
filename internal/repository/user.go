package repository

import (
	"bluebell/internal/entity"
	"gorm.io/gorm"
)

var _ UserRepo = (*userRepo)(nil)

type UserRepo interface {
	IsDuplicateUsername(username string) bool
	IsDuplicateEmail(email string) bool
	InsertUser(user *entity.User) error
}

type userRepo struct {
	db *gorm.DB
}

func (r *userRepo) InsertUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) IsDuplicateUsername(username string) bool {
	return r.db.Where("username = ?", username).First(&entity.User{}).Error == nil
}

func (r *userRepo) IsDuplicateEmail(email string) bool {
	return r.db.Where("email = ?", email).First(&entity.User{}).Error == nil
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}
