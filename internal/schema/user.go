package schema

type UserSignupRequest struct {
	Username string `json:"username" binding:"required,gte=2,lte=24"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanum,gte=6,lte=16"`
}
