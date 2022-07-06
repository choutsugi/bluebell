package codes

const defaultMessage = "CODE_INTERNAL_SERVER_ERROR"

const (
	Success                  = "CODE_SUCCESS"
	ErrInternalServerError   = "CODE_INTERNAL_SERVER_ERROR"
	ErrRequestParamsInvalid  = "CODE_REQUEST_PARAMS_INVALID"
	ErrEmailInvalid          = "CODE_EMAIL_INVALID"
	ErrEmailHasRegistered    = "CODE_EMAIL_HAS_REGISTERED"
	ErrUsernameHasRegistered = "CODE_USERNAME_HAS_REGISTERED"
	ErrPasswordInvalid       = "CODE_PASSWORD_INVALID"
	ErrTokenInvalid          = "CODE_TOKEN_INVALID"
	ErrTokenMissing          = "CODE_TOKEN_MISSING"
	ErrCommunityNotFound     = "CODE_COMMUNITY_NOT_FOUND"
	ErrPostNotFound          = "CODE_POST_NOT_FOUND"
	ErrUserNotCertified      = "CODE_USER_NOT_CERTIFIED"
)

var maps = map[string]string{
	Success:                  "成功",
	ErrInternalServerError:   "服务器异常",
	ErrRequestParamsInvalid:  "请求参数错误",
	ErrEmailInvalid:          "邮箱或密码错误",
	ErrEmailHasRegistered:    "邮箱已注册",
	ErrUsernameHasRegistered: "用户名已注册",
	ErrPasswordInvalid:       "用户名或密码错误",
	ErrTokenInvalid:          "Token无效",
	ErrTokenMissing:          "Token缺失",
	ErrCommunityNotFound:     "社区不存在",
	ErrPostNotFound:          "帖子不存在",
	ErrUserNotCertified:      "用户未登录",
}

func GetMsg(code string) string {
	msg, exist := maps[code]
	if !exist {
		return defaultMessage
	}
	return msg
}
