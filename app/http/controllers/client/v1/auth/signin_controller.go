// Package auth 处理用户身份认证相关逻辑
package auth

import (
	v1 "api/app/http/controllers/client/v1"
	"api/pkg/jwt"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

// SigninController 注册控制器
type SigninController struct {
    v1.BaseAPIController
}

func (sc *SigninController) SignIn(c *gin.Context) {

	// user := user.Get("1")	
	

	token := jwt.NewJWT().IssueToken("11", "xiaolu", "user", "shop-user")

        response.JSON(c, gin.H{
            "token": token,
        })

	// response.Data(c, gin.H{"login": "login successful"})
}

func (sc *SigninController) RefreshToken(c *gin.Context) {

	token,err := jwt.NewJWT().RefreshToken(c)
	if err != nil {
		response.Error(c,err,"token refresh failed")
	} else {
		response.JSON(c, gin.H{
			"token": token,
		})
	}
}
