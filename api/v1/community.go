package v1

import (
	"bluebell/internal/pkg/errx"
	"bluebell/internal/pkg/result"
	"bluebell/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

var _ CommunityApi = (*communityApi)(nil)

type CommunityApi interface {
	FetchAll(ctx *gin.Context)
	FetchOneById(ctx *gin.Context)
}

type communityApi struct {
	service service.CommunityService
}

func (api *communityApi) FetchOneById(ctx *gin.Context) {
	cidStr := ctx.Param("id")
	cid, err := strconv.ParseInt(cidStr, 10, 64)
	if err != nil {
		result.Error(ctx, errx.ErrRequestParamsInvalid)
		return
	}

	community, err := api.service.FetchOneById(cid)
	if err != nil {
		result.Error(ctx, err)
		return
	}
	result.Success(ctx, community)
}

func (api *communityApi) FetchAll(ctx *gin.Context) {
	communities, err := api.service.FetchAll()
	if err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, communities)
}

func newCommunityApi(communityService service.CommunityService) CommunityApi {
	return &communityApi{communityService}
}
