package cache

import (
	"bluebell/internal/conf"
	"bluebell/internal/consts"
	"bluebell/internal/pkg/errx"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

var _ VoteCache = (*voteCache)(nil)

type VoteCache interface {
	Insert(id, communityID string) (err error)
	Vote(id, uid string, opinion float64) (err error)
	FetchIDsWithOrder(start, end int64, orderBy string) ([]string, error)
	FetchIDsByCommunityWithOrder(communityID, start, end int64, orderBy string) ([]string, error)
	CountLikes(ids []string) (data []int64, err error)
}

type voteCache struct {
	rdb  *redis.Client
	conf *conf.Ranking
}

func (cache *voteCache) FetchIDsByCommunityWithOrder(communityID, start, end int64, orderBy string) ([]string, error) {

	//社区key
	communityKey := cache.conf.PostCommunityKeyPrefix + strconv.FormatInt(communityID, 10)

	var orderKey string
	switch orderBy {
	case consts.PostOrderByScore:
		orderKey = cache.conf.PostScoreKey
	case consts.PostOrderByTime:
		orderKey = cache.conf.PostTimeKey
	default:
		orderKey = cache.conf.PostTimeKey
	}

	//读取缓存，不存在则先创建
	key := orderKey + strconv.FormatInt(communityID, 10)
	if cache.rdb.Exists(key).Val() < 1 {
		pipeline := cache.rdb.Pipeline()
		//合并Community(set)与PostTime/PostScore(zset)
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, communityKey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	return cache.rdb.ZRevRange(key, start, end).Result()

}

func (cache *voteCache) CountLikes(ids []string) (data []int64, err error) {

	pipeline := cache.rdb.Pipeline()

	for _, id := range ids {
		pipeline.ZCount(cache.conf.PostVotedKeyPrefix+id, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}

	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		val := cmder.(*redis.IntCmd).Val()
		data = append(data, val)
	}

	return
}

func (cache *voteCache) FetchIDsWithOrder(start, end int64, orderBy string) ([]string, error) {

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

	votedKey := cache.conf.PostVotedKeyPrefix + id

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
	if opinion == originOpinion {
		return errx.ErrDuplicateVotingNotAllowed
	} else if opinion > originOpinion {
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

func (cache *voteCache) Insert(id, communityID string) (err error) {

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

	communityKey := cache.conf.PostCommunityKeyPrefix + communityID
	pipeline.SAdd(communityKey, id)

	_, err = pipeline.Exec()
	return err
}

func NewVoteCache(rdb *redis.Client, conf *conf.Ranking) VoteCache {
	return &voteCache{rdb: rdb, conf: conf}
}
