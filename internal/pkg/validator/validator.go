package validator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

type Err = validator.ValidationErrors

var (
	ZhTrans = genTrans("zh")
	EnTrans = genTrans("en")
)

func GetTranslator(ctx *gin.Context) ut.Translator {
	acceptLanguage := ctx.GetHeader("Accept-Language")
	if strings.Contains(acceptLanguage, "zh") {
		return ZhTrans
	} else {
		return EnTrans
	}
}

func genTrans(locale string) ut.Translator {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT, zhT, enT)

		translator, ok := uni.GetTranslator(locale)
		if !ok {
			panic(fmt.Sprintf("uni.GetTranslator(%s) failed", locale))
		}

		var err error
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(validate, translator)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(validate, translator)
		default:
			err = enTranslations.RegisterDefaultTranslations(validate, translator)
		}
		if err != nil {
			panic(err)
		}

		validate.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		return translator
	}
	panic("validator's engine of gin is not validator.Validate")
}

func TrimFieldPrefix(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
