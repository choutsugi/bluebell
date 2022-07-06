package schema

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

type PostFetchListRequest struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}
