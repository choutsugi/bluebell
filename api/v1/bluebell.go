package v1

import "bluebell/internal/service"

type Api struct {
	User      UserApi
	Community CommunityApi
}

func Register(userService service.UserService, communityService service.CommunityService) Api {
	return Api{
		User:      newUserApi(userService),
		Community: newCommunityApi(communityService),
	}
}
