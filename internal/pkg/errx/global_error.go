package errx

import "bluebell/internal/pkg/codes"

var (
	ErrInternalServerError       = New(codes.ErrInternalServerError)
	ErrInternalServerBusy        = New(codes.ErrInternalServerBusy)
	ErrRequestParamsInvalid      = New(codes.ErrRequestParamsInvalid)
	ErrEmailInvalid              = New(codes.ErrEmailInvalid)
	ErrEmailHasRegistered        = New(codes.ErrEmailHasRegistered)
	ErrUsernameHasRegistered     = New(codes.ErrUsernameHasRegistered)
	ErrPasswordInvalid           = New(codes.ErrPasswordInvalid)
	ErrTokenMissing              = New(codes.ErrTokenMissing)
	ErrTokenInvalid              = New(codes.ErrTokenInvalid)
	ErrCommunityNotFound         = New(codes.ErrCommunityNotFound)
	ErrPostNotFound              = New(codes.ErrPostNotFound)
	ErrUserNotCertified          = New(codes.ErrUserNotCertified)
	ErrBeyondVotingPeriod        = New(codes.ErrBeyondVotingPeriod)
	ErrDuplicateVotingNotAllowed = New(codes.ErrDuplicateVotingNotAllowed)
)
