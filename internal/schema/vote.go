package schema

type VotePostRequest struct {
	PostID  int64 `json:"post_id,string" binding:"required"`     //帖子ID
	Opinion int   `json:"opinion,string" binding:"oneof=1 0 -1"` //观点：赞成（1）、反对（-1）、取消（0）
}
