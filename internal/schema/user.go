package schema

type UserSignupRequest struct {
	Username string `json:"username" binding:"required,gte=2,lte=24"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,gte=6,lte=16"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,gte=6,lte=16"`
}

type UserLoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
