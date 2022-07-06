package v1

import (
	"bluebell/internal/pkg/errx"
	"bluebell/internal/pkg/result"
	"bluebell/internal/schema"
	"bluebell/internal/service"
	"github.com/gin-gonic/gin"
)

var _ VoteApi = (*voteApi)(nil)

type VoteApi interface {
	PostVote(ctx *gin.Context)
}

type voteApi struct {
	service service.VoteService
}

func (api *voteApi) PostVote(ctx *gin.Context) {
	req := new(schema.VotePostRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		result.Error(ctx, err)
		return
	}

	uidData, exists := ctx.Get("uid")
	if !exists {
		result.Error(ctx, errx.ErrUserNotCertified)
		return
	}
	uid := uidData.(int64)

	if err := api.service.Vote(uid, req); err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, nil)
}

func newVoteApi(service service.VoteService) VoteApi {
	return &voteApi{service: service}
}
