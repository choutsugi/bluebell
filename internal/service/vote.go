package service

import (
	"bluebell/internal/cache"
	"bluebell/internal/conf"
	"bluebell/internal/pkg/errx"
	"bluebell/internal/schema"
	"math"
	"strconv"
	"time"
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
	postVotedKey := s.conf.PostVotedPrefix + strconv.FormatInt(req.PostID, 10)

	//获取帖子发布时间
	postTime, err := s.cache.GetPostTime(s.conf.PostTimeKey, strconv.FormatInt(req.PostID, 10))
	if err != nil {
		return errx.ErrBeyondVotingPeriod
	}

	//帖子发布后一周内允许投票
	if float64(time.Now().Unix())-postTime > float64(s.conf.PostVotingPeriod) {
		return errx.ErrBeyondVotingPeriod
	}

	//查询用户投票记录
	originOpinion, err := s.cache.GetVoteRecord(postVotedKey, uidStr)
	if err != nil {
		return err
	}

	//计算分数
	var orientation float64
	if opinion > originOpinion {
		orientation = 1
	} else {
		orientation = -1
	}
	diff := math.Abs(opinion - originOpinion)
	score := orientation * diff * s.conf.PostVoteUnitScore

	//更新投票排行榜分数
	if err = s.cache.UpdateRankingScore(s.conf.PostScoreKey, postIdStr, score); err != nil {
		return err
	}

	if opinion == 0 {
		//删除用户投票信息
		if err = s.cache.DeleteVoteRecord(postVotedKey, uidStr); err != nil {
			return err
		}
	} else {
		//记录用户投票信息
		if err = s.cache.InsertVoteRecord(postVotedKey, uidStr, opinion); err != nil {
			return err
		}
	}

	return nil
}

func NewVoteService(cache cache.VoteCache, conf *conf.Ranking) VoteService {
	return &voteService{cache: cache, conf: conf}
}
