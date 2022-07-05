package v1

import (
	"bluebell/internal/pkg/result"
	"bluebell/internal/service"
	"github.com/gin-gonic/gin"
)

var _ CommunityApi = (*communityApi)(nil)

type CommunityApi interface {
	FetchAll(ctx *gin.Context)
}

type communityApi struct {
	service service.CommunityService
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
