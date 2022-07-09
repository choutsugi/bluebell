package service

import (
	"bluebell/internal/conf"
	"bluebell/internal/data/cache"
	"bluebell/internal/schema"
	"strconv"
)

var _ VoteService = (*voteService)(nil)

type VoteService interface {
	Vote(userID int64, req *schema.VotePostRequest) (err error)
}

type voteService struct {
	cache cache.VoteCache
	conf  *conf.Ranking
}

func (s *voteService) Vote(userID int64, req *schema.VotePostRequest) (err error) {

	id := strconv.FormatInt(req.PostID, 10)
	uid := strconv.FormatInt(userID, 10)
	opinion := float64(req.Opinion)

	return s.cache.Vote(id, uid, opinion)
}

func NewVoteService(cache cache.VoteCache, conf *conf.Ranking) VoteService {
	return &voteService{cache: cache, conf: conf}
}
