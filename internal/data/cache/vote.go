package cache

import (
	"bluebell/internal/conf"
	"bluebell/internal/consts"
	"bluebell/internal/pkg/errx"
	"github.com/go-redis/redis"
	"math"
	"time"
)

var _ VoteCache = (*voteCache)(nil)

type VoteCache interface {
	Insert(id string) (err error)
	Vote(id, uid string, opinion float64) (err error)
	FetchIDs(start, end int64, orderBy string) ([]string, error)
}

type voteCache struct {
	rdb  *redis.Client
	conf *conf.Ranking
}

func (cache *voteCache) FetchIDs(start, end int64, orderBy string) ([]string, error) {

	var key string
	switch orderBy {
	case consts.PostOrderByScore:
		key = cache.conf.PostScoreKey
	case consts.PostOrderByTime:
		key = cache.conf.PostTimeKey
	default:
		key = cache.conf.PostTimeKey
	}
	//倒序
	return cache.rdb.ZRevRange(key, start, end).Result()
}

func (cache *voteCache) Vote(id, uid string, opinion float64) (err error) {

	votedKey := cache.conf.PostVotedPrefix + id

	//获取帖子发布时间
	postTime, err := cache.rdb.ZScore(cache.conf.PostTimeKey, id).Result()
	if err != nil {
		return errx.ErrBeyondVotingPeriod
	}

	//帖子发布后一周内允许投票
	if float64(time.Now().Unix())-postTime > cache.conf.PostVotingPeriod {
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
	score := orientation * diff * cache.conf.PostVoteUnitScore

	//开始事务
	pipeline := cache.rdb.TxPipeline()

	//更新投票排行榜分数
	pipeline.ZIncrBy(cache.conf.PostScoreKey, score, id)

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

func (cache *voteCache) Insert(id string) (err error) {

	//开启事务
	pipeline := cache.rdb.TxPipeline()

	pipeline.ZAdd(cache.conf.PostTimeKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: id,
	})

	pipeline.ZAdd(cache.conf.PostScoreKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: id,
	})

	_, err = pipeline.Exec()
	return err
}

func NewVoteCache(rdb *redis.Client, conf *conf.Ranking) VoteCache {
	return &voteCache{rdb: rdb, conf: conf}
}
