package repository

import (
	"bluebell/internal/entity"
	"bluebell/internal/pkg/errx"
	"gorm.io/gorm"
)

var _ CommunityRepo = (*communityRepo)(nil)

type CommunityRepo interface {
	FetchAll() (communities []*entity.Community, err error)
	FetchOneById(cid int64) (community *entity.Community, err error)
}

type communityRepo struct {
	db *gorm.DB
}

func (repo *communityRepo) FetchOneById(cid int64) (community *entity.Community, err error) {
	community = new(entity.Community)
	if err = repo.db.Where("cid = ?", cid).Take(community).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errx.ErrCommunityNotFound
		}
		return nil, err
	}
	return
}

func (repo *communityRepo) FetchAll() (communities []*entity.Community, err error) {
	communities = make([]*entity.Community, 0)
	err = repo.db.Find(&communities).Error
	return
}

func NewCommunityRepo(db *gorm.DB) CommunityRepo {
	return &communityRepo{db: db}
}
