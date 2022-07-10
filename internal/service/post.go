package service

import (
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
	FetchListWithOrder(req *schema.PostFetchWithOrderRequest) (resp []*schema.PostDetail, err error)
	FetchListByPaginate(req *schema.PostFetchPaginateRequest) (posts []*entity.Post, err error)
	FetchListByCommunityWithOrder(req *schema.PostFetchByCommunityWithOrderRequest) (resp []*schema.PostDetail, err error)
}

type postService struct {
	postRepo      repo.PostRepo
	userRepo      repo.UserRepo
	communityRepo repo.CommunityRepo
	cache         cache.VoteCache
}

func (s *postService) FetchListByCommunityWithOrder(req *schema.PostFetchByCommunityWithOrderRequest) (resp []*schema.PostDetail, err error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}

	//计算redis数据起始点
	start := int64((req.PageNum - 1) * req.PageSize)
	end := start + int64(req.PageSize) - 1

	//获取id
	idsStr, err := s.cache.FetchIDsByCommunityWithOrder(req.CommunityID, start, end, req.OrderBy)
	if err != nil {
		return nil, err
	}

	likes, err := s.cache.CountLikes(idsStr)
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

	posts, err := s.postRepo.FetchListByIDs(ids)
	if err != nil {
		return nil, err
	}

	resp = make([]*schema.PostDetail, 0, len(posts))

	for _, post := range posts {
		user, err := s.userRepo.FetchUserByID(post.ID)
		if err != nil {
			return nil, err
		}
		community, err := s.communityRepo.FetchOneById(post.CommunityId)
		if err != nil {
			return nil, err
		}
		postDetail := &schema.PostDetail{
			AuthorName: user.Username,
			Community:  community,
			Likes:      likes[post.ID],
			Post:       post,
		}

		resp = append(resp, postDetail)
	}

	return
}

func (s *postService) FetchListByPaginate(req *schema.PostFetchPaginateRequest) (posts []*entity.Post, err error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	offset := (req.PageNum - 1) * req.PageSize
	limit := req.PageSize
	return s.postRepo.FetchListByPaginate(offset, limit)
}

func (s *postService) FetchListWithOrder(req *schema.PostFetchWithOrderRequest) (resp []*schema.PostDetail, err error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}

	//计算redis数据起始点
	start := int64((req.PageNum - 1) * req.PageSize)
	end := start + int64(req.PageSize) - 1

	//获取id
	idsStr, err := s.cache.FetchIDsWithOrder(start, end, req.OrderBy)
	if err != nil {
		return nil, err
	}

	likes, err := s.cache.CountLikes(idsStr)
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

	posts, err := s.postRepo.FetchListByIDs(ids)
	if err != nil {
		return nil, err
	}

	resp = make([]*schema.PostDetail, 0, len(posts))

	for _, post := range posts {
		user, err := s.userRepo.FetchUserByID(post.ID)
		if err != nil {
			return nil, err
		}
		community, err := s.communityRepo.FetchOneById(post.CommunityId)
		if err != nil {
			return nil, err
		}
		postDetail := &schema.PostDetail{
			AuthorName: user.Username,
			Community:  community,
			Likes:      likes[post.ID],
			Post:       post,
		}

		resp = append(resp, postDetail)
	}

	return
}

func (s *postService) Create(req *schema.PostCreateRequest) (err error) {

	post := entity.Post{
		ID:          snowflake.GenerateID(),
		Title:       req.Title,
		Content:     req.Content,
		AuthorId:    req.AuthorId,
		CommunityId: req.CommunityId,
	}
	if err = s.postRepo.Insert(&post); err != nil {
		return err
	}
	id := strconv.FormatInt(post.ID, 10)
	communityID := strconv.FormatInt(post.CommunityId, 10)

	return s.cache.Insert(id, communityID)
}

func (s *postService) Delete(postId int64) (err error) {
	post := entity.Post{
		ID: postId,
	}
	return s.postRepo.Delete(&post)
}

func (s *postService) Update(req *schema.PostUpdateRequest) (err error) {
	post := entity.Post{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
	}
	return s.postRepo.Update(&post)
}

func (s *postService) FetchByID(postId int64) (post *entity.Post, err error) {
	return s.postRepo.FetchByID(postId)
}

func (s *postService) FetchAll() (posts []*entity.Post, err error) {
	return s.postRepo.FetchAll()
}

func NewPostService(postRepo repo.PostRepo, userRepo repo.UserRepo, communityRepo repo.CommunityRepo, cache cache.VoteCache) PostService {
	return &postService{
		postRepo:      postRepo,
		userRepo:      userRepo,
		communityRepo: communityRepo,
		cache:         cache,
	}
}
