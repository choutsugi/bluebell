package v1

import (
	"bluebell/internal/pkg/errx"
	"bluebell/internal/pkg/result"
	"bluebell/internal/schema"
	"bluebell/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

var _ PostApi = (*postApi)(nil)

type PostApi interface {
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	FetchById(ctx *gin.Context)
	FetchAll(ctx *gin.Context)
	FetchListWithOrder(ctx *gin.Context)
	FetchListByPaginate(ctx *gin.Context)
}

type postApi struct {
	service service.PostService
}

func (api *postApi) FetchListByPaginate(ctx *gin.Context) {
	req := &schema.PostFetchPaginateRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		result.Error(ctx, err)
		return
	}

	posts, err := api.service.FetchListByPaginate(req)
	if err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, posts)
}

func (api *postApi) FetchListWithOrder(ctx *gin.Context) {
	//可指定默认参数
	req := &schema.PostFetch{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		result.Error(ctx, err)
		return
	}

	var (
		posts []*schema.PostDetail
		err   error
	)
	if req.CommunityID == 0 {
		posts, err = api.service.FetchListWithOrder(req)
	} else {
		posts, err = api.service.FetchListByCommunityWithOrder(req)

	}
	if err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, posts)
}

func (api *postApi) Create(ctx *gin.Context) {
	req := new(schema.PostCreateRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		result.Error(ctx, err)
		return
	}
	uid, exists := ctx.Get("uid")
	if !exists {
		result.Error(ctx, errx.ErrUserNotCertified)
		return
	}
	req.AuthorId = uid.(int64)

	if err := api.service.Create(req); err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, nil)
}

func (api *postApi) Delete(ctx *gin.Context) {
	pidStr := ctx.Param("id")
	postId, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		result.Error(ctx, errx.ErrRequestParamsInvalid)
		return
	}

	if err = api.service.Delete(postId); err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, nil)
}

func (api *postApi) Update(ctx *gin.Context) {
	req := new(schema.PostUpdateRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		result.Error(ctx, err)
		return
	}

	if err := api.service.Update(req); err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, nil)
}

func (api *postApi) FetchById(ctx *gin.Context) {
	pidStr := ctx.Param("id")
	postId, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		result.Error(ctx, errx.ErrRequestParamsInvalid)
		return
	}

	post, err := api.service.FetchByID(postId)
	if err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, post)
}

func (api *postApi) FetchAll(ctx *gin.Context) {
	posts, err := api.service.FetchAll()
	if err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, posts)
}

func newPostApi(postService service.PostService) PostApi {
	return &postApi{postService}
}
