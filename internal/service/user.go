package service

import (
	"bluebell/internal/entity"
	"bluebell/internal/pkg/auth"
	"bluebell/internal/pkg/errx"
	"bluebell/internal/repository"
	"bluebell/internal/schema"
	"bluebell/pkg/encrypt"
	"bluebell/pkg/snowflake"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	Signup(req *schema.UserSignupRequest) (err error)
	Login(req *schema.UserLoginRequest) (resp *schema.UserLoginResponse, err error)
}

type userService struct {
	repo repository.UserRepo
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
		Uid:      snowflake.GenerateID(),
		Username: req.Username,
		Email:    req.Email,
	}
	user.Password, user.Salt, err = encrypt.Encrypt(req.Password)
	if err != nil {
		return err
	}

	return s.repo.InsertUser(&user)
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}
