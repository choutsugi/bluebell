package service

import (
	"bluebell/internal/conf"
	"bluebell/internal/data/cache"
	"bluebell/internal/data/repo"
	"bluebell/internal/entity"
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
	FetchListByOrder(req *schema.PostFetchOrderRequest) (posts []*entity.Post, err error)
	FetchListByPaginate(req *schema.PostFetchPaginateRequest) (posts []*entity.Post, err error)
}

type postService struct {
	repo  repo.PostRepo
	cache cache.VoteCache
	conf  *conf.Ranking
}

func (s *postService) FetchListByPaginate(req *schema.PostFetchPaginateRequest) (posts []*entity.Post, err error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	offset := (req.PageNum - 1) * req.PageSize
	limit := req.PageSize
	return s.repo.FetchListByPaginate(offset, limit)
}

func (s *postService) FetchListByOrder(req *schema.PostFetchOrderRequest) (posts []*entity.Post, err error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}

	//计算redis数据起始点
	start := int64((req.PageNum - 1) * req.PageSize)
	end := start + int64(req.PageSize) - 1

	//获取id
	idsStr, err := s.cache.FetchIDs(start, end, req.OrderBy)
	if err != nil {
		return nil, err
	}

	var ids []int64
	for _, item := range idsStr {
		id, err := strconv.ParseInt(item, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, nil
	}

	return s.repo.FetchListByIDs(ids)
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

	return s.cache.Insert(s.conf.PostTimeKey, s.conf.PostScoreKey, strconv.FormatInt(post.ID, 10))
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

func NewPostService(repo repo.PostRepo, cache cache.VoteCache, conf *conf.Ranking) PostService {
	return &postService{repo: repo, cache: cache, conf: conf}
}
