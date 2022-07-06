package service

import (
	"bluebell/internal/data/repo"
	"bluebell/internal/entity"
)

var _ CommunityService = (*communityService)(nil)

type CommunityService interface {
	FetchAll() (communities []*entity.Community, err error)
	FetchOneById(cid int64) (community *entity.Community, err error)
}

type communityService struct {
	repo repo.CommunityRepo
}

func (s *communityService) FetchOneById(cid int64) (community *entity.Community, err error) {
	return s.repo.FetchOneById(cid)
}

func (s *communityService) FetchAll() (communities []*entity.Community, err error) {
	return s.repo.FetchAll()
}

func NewCommunityService(repo repo.CommunityRepo) CommunityService {
	return &communityService{repo: repo}
}
