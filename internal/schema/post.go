package schema

import "bluebell/internal/entity"

type PostCreateRequest struct {
	Title       string `json:"title" binding:"required,gte=6,lte=128"`
	Content     string `json:"content" binding:"required,gte=6"`
	CommunityId int64  `json:"community_id" binding:"required,numeric"`
	AuthorId    int64
}

type PostUpdateRequest struct {
	ID      int64  `json:"id" binding:"required,numeric"`
	Title   string `json:"title" binding:"required,gte=6,lte=128"`
	Content string `json:"content" binding:"required,gte=6"`
}

type PostFetchPaginateRequest struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

type PostFetchWithOrderRequest struct {
	OrderBy  string `form:"order_by"`
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
}

type PostFetchByCommunityWithOrderRequest struct {
	CommunityID int64  `form:"community_id"`
	OrderBy     string `form:"order_by"`
	PageNum     int    `form:"page_num"`
	PageSize    int    `form:"page_size"`
}

type PostDetail struct {
	AuthorName string
	Community  *entity.Community
	Likes      int64
	Post       *entity.Post
}
