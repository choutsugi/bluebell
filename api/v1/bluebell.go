package v1

import "bluebell/internal/service"

type Api struct {
	User      UserApi
	Community CommunityApi
	Post      PostApi
	Vote      VoteApi
}

func Register(userService service.UserService, communityService service.CommunityService, postService service.PostService, voteService service.VoteService) Api {
	return Api{
		User:      newUserApi(userService),
		Community: newCommunityApi(communityService),
		Post:      newPostApi(postService),
		Vote:      newVoteApi(voteService),
	}
}
