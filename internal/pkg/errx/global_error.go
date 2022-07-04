package errx

import "bluebell/internal/pkg/codes"

var (
	ErrInternalServerError   = New(codes.ErrInternalServerError)
	ErrEmailInvalid          = New(codes.ErrEmailInvalid)
	ErrEmailHasRegistered    = New(codes.ErrEmailHasRegistered)
	ErrUsernameHasRegistered = New(codes.ErrUsernameHasRegistered)
	ErrPasswordInvalid       = New(codes.ErrPasswordInvalid)
	ErrTokenMissing          = New(codes.ErrTokenMissing)
	ErrTokenInvalid          = New(codes.ErrTokenInvalid)
)
