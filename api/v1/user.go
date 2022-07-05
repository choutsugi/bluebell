package v1

import (
	"bluebell/internal/pkg/result"
	"bluebell/internal/schema"
	"bluebell/internal/service"
	"github.com/gin-gonic/gin"
)

var _ UserApi = (*userApi)(nil)

type UserApi interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type userApi struct {
	service service.UserService
}

func (api *userApi) Signup(ctx *gin.Context) {
	var req schema.UserSignupRequest
	if err := ctx.ShouldBind(&req); err != nil {
		result.Error(ctx, err)
		return
	}

	if err := api.service.Signup(&req); err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, nil)
}

func (api *userApi) Login(ctx *gin.Context) {
	var req schema.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		result.Error(ctx, err)
		return
	}

	resp, err := api.service.Login(&req)
	if err != nil {
		result.Error(ctx, err)
		return
	}

	result.Success(ctx, resp)
}

func newUserApi(userService service.UserService) UserApi {
	return &userApi{service: userService}
}
