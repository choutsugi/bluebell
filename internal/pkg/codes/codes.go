package codes

const defaultMessage = "CODE_INTERNAL_SERVER_ERROR"

const (
	Success                  = "CODE_SUCCESS"
	ErrInternalServerError   = "CODE_INTERNAL_SERVER_ERROR"
	ErrRequestParamsInvalid  = "CODE_REQUEST_PARAMS_INVALID"
	ErrUserNotFound          = "CODE_USER_NOT_FOUND"
	ErrEmailHasRegistered    = "CODE_EMAIL_HAS_REGISTERED"
	ErrUsernameHasRegistered = "CODE_USERNAME_HAS_REGISTERED"
)

var maps = map[string]string{
	Success:                  "成功",
	ErrInternalServerError:   "服务器异常",
	ErrRequestParamsInvalid:  "请求参数错误",
	ErrUserNotFound:          "用户名或密码错误",
	ErrEmailHasRegistered:    "邮箱已注册",
	ErrUsernameHasRegistered: "用户名已注册",
}

func GetMsg(code string) string {
	msg, exist := maps[code]
	if !exist {
		return defaultMessage
	}
	return msg
}
