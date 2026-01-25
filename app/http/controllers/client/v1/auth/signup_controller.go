// Package auth 处理用户身份认证相关逻辑
package auth

import (
	v1 "api/app/http/controllers/client/v1"
	"api/app/requests"

	request "api/app/requests"
	authRequest "api/app/requests/client/auth"

	"api/pkg/logger"

	"github.com/gin-gonic/gin"

	"api/pkg/response"
)

// SignupController 注册控制器
type SignupController struct {
    v1.BaseAPIController
}

func (*SignupController) IsEmailExist(c *gin.Context) {
 	var req authRequest.SignupRequest

	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err,
	// 	})
	
	// 	return
	// }
	logger.DebugString("", "","")

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "ShouldBindJSON ERR")
		return
	}

	if err := request.Validate.Struct(&req); err != nil {
		errMsg := requests.MakeErrorMsg(c,err)
		response.ValidationError(c,errMsg)
		return
	}

	response.JSON(c,gin.H{"exists": false})

}

// SignupUsingEmail 使用 Email + 验证码进行注册
func (sc *SignupController) SignupUsingEmail(c *gin.Context) {

	
}
