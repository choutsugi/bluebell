package errx

import "bluebell/internal/pkg/codes"

type Error struct {
	Code    string
	Message string
}

func (err *Error) Error() string {
	return err.Message
}

func New(code string) *Error {
	return &Error{
		Code:    code,
		Message: codes.GetMsg(code),
	}
}

func NewWithData(code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}
