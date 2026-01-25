package requests

import (
	"api/pkg/i18n"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 全局 validator，供 controller / request 使用
var Validate *validator.Validate

// TranslateValidationErrors
// 将 validator 错误转换为 i18n key，不做语言翻译
func TranslateValidationErrors(err error) map[string]string {
	result := make(map[string]string)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			// email.required / phone.phone_not_exists
			key := e.Field() + "." + e.Tag()
			result[e.Field()] = key
		}
	}

	return result
}

func MakeErrorMsg(c *gin.Context, err error) (errors map[string]string) {
	keys := TranslateValidationErrors(err)
	lang, _ := c.Get("lang")
	errors = make(map[string]string)
	for field, key := range keys {
		errors[field] = i18n.T(lang.(string), key)
	}
	return errors
}

