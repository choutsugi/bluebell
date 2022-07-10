package repo

import (
	"bluebell/internal/entity"
	"bluebell/internal/pkg/errx"
	"gorm.io/gorm"
)

var _ UserRepo = (*userRepo)(nil)

type UserRepo interface {
	IsDuplicateUsername(username string) bool
	IsDuplicateEmail(email string) bool
	InsertUser(user *entity.User) error
	FetchUserByEmail(email string) (user *entity.User, err error)
	FetchUserByID(uid int64) (user *entity.User, err error)
}

type userRepo struct {
	db *gorm.DB
}

func (r *userRepo) FetchUserByID(uid int64) (user *entity.User, err error) {
	user = new(entity.User)
	if err = r.db.First(user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errx.ErrEmailInvalid
		}
		return nil, err
	}
	return
}

func (r *userRepo) FetchUserByEmail(email string) (user *entity.User, err error) {
	user = new(entity.User)
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errx.ErrEmailInvalid
		}
		return nil, err
	}
	return
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
