package cache

import (
	"github.com/go-redis/redis"
	"time"
)

var _ VoteCache = (*voteCache)(nil)

type VoteCache interface {
	JoinRanking(postTimeKey, postId string) (err error)
	GetPostTime(postTimeKey, postId string) (postTime float64, err error)
	GetVoteRecord(postVotedKey, uid string) (opinion float64, err error)
	UpdateRankingScore(postScoreKey, postId string, score float64) (err error)
	DeleteVoteRecord(postVotedKey, uid string) (err error)
	InsertVoteRecord(postVotedKey, uid string, opinion float64) (err error)
}

type voteCache struct {
	rdb *redis.Client
}

func (cache *voteCache) InsertVoteRecord(postVotedKey, uid string, opinion float64) (err error) {
	_, err = cache.rdb.ZAdd(postVotedKey, redis.Z{
		Score:  opinion,
		Member: uid,
	}).Result()
	return
}

func (cache *voteCache) DeleteVoteRecord(postVotedKey, uid string) (err error) {
	_, err = cache.rdb.ZRem(postVotedKey, uid).Result()
	return
}

func (cache *voteCache) UpdateRankingScore(postScoreKey, postId string, score float64) (err error) {
	_, err = cache.rdb.ZIncrBy(postScoreKey, score, postId).Result()
	return
}

func (cache *voteCache) GetVoteRecord(postVotedKey, uid string) (opinion float64, err error) {
	return cache.rdb.ZScore(postVotedKey, uid).Result()
}

func (cache *voteCache) GetPostTime(postTimeKey, postId string) (postTime float64, err error) {
	return cache.rdb.ZScore(postTimeKey, postId).Result()
}

func (cache *voteCache) JoinRanking(postTimeKey, postId string) (err error) {
	_, err = cache.rdb.ZAdd(postTimeKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	}).Result()
	return
}

func NewVoteCache(rdb *redis.Client) VoteCache {
	return &voteCache{rdb: rdb}
}
