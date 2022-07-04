package service

import (
	"bluebell/internal/entity"
	"bluebell/internal/pkg/errx"
	"bluebell/internal/repository"
	"bluebell/internal/schema"
	"bluebell/pkg/encrypt"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	Signup(req *schema.UserSignupRequest) (err error)
}

type userService struct {
	repo repository.UserRepo
}

func (s *userService) Signup(req *schema.UserSignupRequest) (err error) {

	if s.repo.IsDuplicateUsername(req.Username) {
		return errx.ErrUsernameHasRegistered
	}

	if s.repo.IsDuplicateEmail(req.Username) {
		return errx.ErrEmailHasRegistered
	}

	user := new(entity.User)
	user.Username = req.Username
	user.Email = req.Email

	user.Password, user.Salt, err = encrypt.Encrypt(req.Password)
	if err != nil {
		return err
	}

	return s.repo.InsertUser(user)
}

func NewUserService(repo repository.UserRepo) UserService {
	return &userService{repo: repo}
}
