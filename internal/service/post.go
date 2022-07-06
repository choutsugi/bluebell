package service

import (
	"bluebell/internal/cache"
	"bluebell/internal/conf"
	"bluebell/internal/entity"
	"bluebell/internal/repository"
	"bluebell/internal/schema"
	"bluebell/pkg/snowflake"
	"strconv"
)

var _ PostService = (*postService)(nil)

type PostService interface {
	Create(req *schema.PostCreateRequest) (err error)
	Delete(postId int64) (err error)
	Update(req *schema.PostUpdateRequest) (err error)
	FetchByID(postId int64) (post *entity.Post, err error)
	FetchAll() (posts []*entity.Post, err error)
	FetchList(req *schema.PostFetchListRequest) (posts []*entity.Post, err error)
}

type postService struct {
	repo  repository.PostRepo
	cache cache.VoteCache
	conf  *conf.Ranking
}

func (s *postService) FetchList(req *schema.PostFetchListRequest) (posts []*entity.Post, err error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	offset := (req.PageNum - 1) * req.PageSize
	limit := req.PageSize

	return s.repo.FetchList(offset, limit)
}

func (s *postService) Create(req *schema.PostCreateRequest) (err error) {

	post := entity.Post{
		ID:          snowflake.GenerateID(),
		Title:       req.Title,
		Content:     req.Content,
		AuthorId:    req.AuthorId,
		CommunityId: req.CommunityId,
	}
	if err = s.repo.Insert(&post); err != nil {
		return err
	}

	if err = s.cache.JoinRanking(s.conf.PostTimeKey, strconv.FormatInt(post.ID, 10)); err != nil {
		return err
	}

	return nil
}

func (s *postService) Delete(postId int64) (err error) {
	post := entity.Post{
		ID: postId,
	}
	return s.repo.Delete(&post)
}

func (s *postService) Update(req *schema.PostUpdateRequest) (err error) {
	post := entity.Post{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
	}
	return s.repo.Update(&post)
}

func (s *postService) FetchByID(postId int64) (post *entity.Post, err error) {
	return s.repo.FetchByID(postId)
}

func (s *postService) FetchAll() (posts []*entity.Post, err error) {
	return s.repo.FetchAll()
}

func NewPostService(repo repository.PostRepo, cache cache.VoteCache, conf *conf.Ranking) PostService {
	return &postService{repo: repo, cache: cache, conf: conf}
}
