package v1

import "bluebell/internal/service"

type Api struct {
	User UserApi
}

func Register(userService service.UserService) Api {
	return Api{
		User: newUserApi(userService),
	}
}
