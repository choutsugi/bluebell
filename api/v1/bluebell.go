package v1

import "bluebell/internal/service"

type Api struct {
	User      UserApi
	Community CommunityApi
	Post      PostApi
}

func Register(userService service.UserService, communityService service.CommunityService, postService service.PostService) Api {
	return Api{
		User:      newUserApi(userService),
		Community: newCommunityApi(communityService),
		Post:      newPostApi(postService),
	}
}
