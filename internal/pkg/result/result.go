package result

import (
	"bluebell/internal/pkg/codes"
	"bluebell/internal/pkg/errx"
	"bluebell/internal/pkg/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

type result struct {
	Status  bool        `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Nil struct{}

func Success(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = Nil{}
	}

	ctx.JSON(http.StatusOK, &result{
		Status:  true,
		Code:    "CODE_REQUEST_SUCCEEDED",
		Message: "success",
		Data:    data,
	})
}

func Error(ctx *gin.Context, err error) {

	if e, ok := err.(*errx.Error); ok {
		ctx.JSON(http.StatusOK, &result{
			Status:  false,
			Code:    e.Code,
			Message: e.Message,
			Data:    Nil{},
		})
	} else if e, ok := err.(validator.Err); ok {
		ctx.JSON(http.StatusOK, &result{
			Status:  false,
			Code:    codes.ErrRequestParamsInvalid,
			Message: codes.GetMsg(codes.ErrRequestParamsInvalid),
			Data:    validator.TrimFieldPrefix(e.Translate(validator.GetTranslator(ctx))),
		})
	} else {
		ctx.JSON(http.StatusOK, &result{
			Status:  false,
			Code:    codes.ErrInternalServerError,
			Message: codes.GetMsg(codes.ErrInternalServerError),
			Data:    Nil{},
		})
	}
}
