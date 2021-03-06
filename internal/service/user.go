package service

import (
	"bluebell/internal/data/repo"
	"bluebell/internal/entity"
	"bluebell/internal/pkg/auth"
	"bluebell/internal/pkg/errx"
	"bluebell/internal/schema"
	"bluebell/pkg/encrypt"
	"github.com/golang-jwt/jwt/v4"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	Signup(req *schema.UserSignupRequest) (err error)
	Login(req *schema.UserLoginRequest) (resp *schema.UserLoginResponse, err error)
	Logout(token *jwt.Token) (err error)
}

type userService struct {
	repo repo.UserRepo
}

func (s *userService) Logout(token *jwt.Token) (err error) {
	return auth.JoinBlacklist(token)
}

func (s *userService) Login(req *schema.UserLoginRequest) (resp *schema.UserLoginResponse, err error) {
	user, err := s.repo.FetchUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !encrypt.Verify(req.Password, user.Password, user.Salt) {
		return nil, errx.ErrPasswordInvalid
	}

	//生成token
	token, err := auth.GenerateToken(user.Uid)
	if err != nil {
		return nil, err
	}

	return &schema.UserLoginResponse{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
		TokenType:   token.TokenType,
	}, nil
}

func (s *userService) Signup(req *schema.UserSignupRequest) (err error) {

	if s.repo.IsDuplicateUsername(req.Username) {
		return errx.ErrUsernameHasRegistered
	}

	if s.repo.IsDuplicateEmail(req.Username) {
		return errx.ErrEmailHasRegistered
	}

	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
	}
	user.Password, user.Salt, err = encrypt.Encrypt(req.Password)
	if err != nil {
		return err
	}

	return s.repo.InsertUser(&user)
}

func NewUserService(repo repo.UserRepo) UserService {
	return &userService{repo: repo}
}
