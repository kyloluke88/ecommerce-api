package routes

import (
	clientAuth "api/app/http/controllers/client/v1/auth"
	"api/pkg/config"

	"api/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterApiRoutes 注册网页相关路由
func RegisterApiRoutes(r *gin.Engine) {

	// 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
	var v1 *gin.RouterGroup

	if len(config.Get[string]("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}
	v1.Use(middlewares.Common())
	
	client := v1.Group("client")
	{
		auth := client.Group("/auth")
		{
			suc := new(clientAuth.SignupController)
			auth.POST("/signup/email/exist", suc.IsEmailExist)
		}
	}

	// admin := v1.Group("admin")
	// {
	// 	auth := admin.Group("/auth")
	// 	{
	// 		suc := new(clientAuth.SignupController)
	// 		auth.POST("/signup/email/exist", suc.SignupUsingEmail)
	// 	}
	// }



	
}
