package repository

import (
	"bluebell/internal/entity"
	"gorm.io/gorm"
)

var _ CommunityRepo = (*communityRepo)(nil)

type CommunityRepo interface {
	FetchAll() (communities []*entity.Community, err error)
}

type communityRepo struct {
	db *gorm.DB
}

func (repo *communityRepo) FetchAll() (communities []*entity.Community, err error) {
	communities = make([]*entity.Community, 0)
	err = repo.db.Find(&communities).Error
	return
}

func NewCommunityRepo(db *gorm.DB) CommunityRepo {
	return &communityRepo{db: db}
}
