package service

import (
	"bluebell/internal/conf"
	"bluebell/internal/data/cache"
	"bluebell/internal/schema"
	"strconv"
)

var _ VoteService = (*voteService)(nil)

type VoteService interface {
	Vote(uid int64, req *schema.VotePostRequest) (err error)
}

type voteService struct {
	cache cache.VoteCache
	conf  *conf.Ranking
}

func (s *voteService) Vote(uid int64, req *schema.VotePostRequest) (err error) {

	uidStr := strconv.FormatInt(uid, 10)
	postIdStr := strconv.FormatInt(req.PostID, 10)
	opinion := float64(req.Opinion)
	period := float64(s.conf.PostVotingPeriod)
	votedKey := s.conf.PostVotedPrefix + strconv.FormatInt(req.PostID, 10)

	return s.cache.Vote(s.conf.PostTimeKey, s.conf.PostScoreKey, votedKey, postIdStr, uidStr, period, s.conf.PostVoteUnitScore, opinion)
}

func NewVoteService(cache cache.VoteCache, conf *conf.Ranking) VoteService {
	return &voteService{cache: cache, conf: conf}
}
