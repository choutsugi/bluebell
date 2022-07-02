package v1

import (
	"bluebell/logic"
	"bluebell/models"
	"bluebell/pkg/validator"
	"github.com/gin-gonic/gin"
	paramValidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

func RegisterHandler(ctx *gin.Context) {

	//1.参数校验
	var param models.ParamRegister
	if err := ctx.ShouldBindJSON(&param); err != nil {
		zap.L().Error("register with invalid param", zap.Error(err))
		errs, ok := err.(paramValidator.ValidationErrors)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msg": validator.RemoveTopStruct(errs.Translate(validator.Trans)),
		})
		return
	}

	//2.注册用户
	if err := logic.Register(param); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}

	//3.返回结果
	ctx.JSON(http.StatusOK, "注册成功")
}
