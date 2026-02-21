// Package auth 处理用户身份认证相关逻辑
package auth

import (
	"api/app/http/controllers/client"
	"api/app/requests"

	request "api/app/requests"
	authRequest "api/app/requests/client/auth"

	"github.com/gin-gonic/gin"

	"api/app/models/user"
	"api/pkg/jwt"
	"api/pkg/response"
)

// SignupController 注册控制器
type SignupController struct {
	client.BaseAPIController
}

func (*SignupController) IsEmailExist(c *gin.Context) {
	var req authRequest.EmailIsExistRequest

	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, err, "ShouldBindJSON ERR")
		return
	}

	if err := request.Validate.Struct(&req); err != nil {
		errMsg := requests.MakeErrorMsg(c, err)
		response.ValidationError(c, errMsg)
		return
	}

	response.JSON(c, gin.H{"exists": false})
}

// SignupUsingEmail 使用 Email + 验证码进行注册
func (sc *SignupController) SignupUsingEmail(c *gin.Context) {

	var req authRequest.SignupEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err, "ShouldBindJSON ERR")
		return
	}

	if err := request.Validate.Struct(&req); err != nil {
		errMsg := requests.MakeErrorMsg(c, err)
		response.ValidationError(c, errMsg)
		return
	}

	userModel := user.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		UserProfile: user.UserProfile{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	}

	userModel.Create()

	if userModel.ID > 0 {
		token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.Email, "user", "shop-user")

		response.JSON(c, gin.H{
			"token": token,
			"data":  userModel,
		})
	} else {
		response.Abort500(c, "user create failed, please try again later")
	}
}

func (sc *SignupController) SignupUseingGoogle(c *gin.Context) {

}
