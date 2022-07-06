package cache

import (
	"bluebell/internal/pkg/errx"
	"github.com/go-redis/redis"
	"math"
	"time"
)

var _ VoteCache = (*voteCache)(nil)

type VoteCache interface {
	InsertPost(postTimeKey, postScoreKey, postId string) (err error)
	VotePost(timeKey, scoreKey, votedKey, postId, uid string, period, unitScore, opinion float64) (err error)
}

type voteCache struct {
	rdb *redis.Client
}

func (cache *voteCache) VotePost(timeKey, scoreKey, votedKey, postId, uid string, period, unitScore, opinion float64) (err error) {
	//获取帖子发布时间
	postTime, err := cache.rdb.ZScore(timeKey, postId).Result()
	if err != nil {
		return errx.ErrBeyondVotingPeriod
	}

	//帖子发布后一周内允许投票
	if float64(time.Now().Unix())-postTime > period {
		return errx.ErrBeyondVotingPeriod
	}

	//查询用户投票记录
	originOpinion, err := cache.rdb.ZScore(votedKey, uid).Result()
	if err != nil && err != redis.Nil {
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
	score := orientation * diff * unitScore

	//开始事务
	pipeline := cache.rdb.TxPipeline()

	//更新投票排行榜分数
	pipeline.ZIncrBy(scoreKey, score, postId)

	if opinion == 0 {
		//删除用户投票信息
		pipeline.ZRem(votedKey, uid)
	} else {
		//记录用户投票信息
		pipeline.ZAdd(votedKey, redis.Z{
			Score:  opinion,
			Member: uid,
		})
	}

	_, err = pipeline.Exec()

	return err
}

func (cache *voteCache) InsertPost(postTimeKey, postScoreKey, postId string) (err error) {

	//开启事务
	pipeline := cache.rdb.TxPipeline()

	pipeline.ZAdd(postTimeKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	pipeline.ZAdd(postScoreKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	_, err = pipeline.Exec()
	return err
}

func NewVoteCache(rdb *redis.Client) VoteCache {
	return &voteCache{rdb: rdb}
}
