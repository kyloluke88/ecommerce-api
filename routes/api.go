package routes

import (
	"api/app/http/controllers/client"
	clientAuth "api/app/http/controllers/client/auth"
	"api/pkg/config"

	"api/app/http/middlewares"

	"github.com/gin-gonic/gin"

	"api/app/http/controllers/admin"
)

// RegisterApiRoutes 注册网页相关路由
func RegisterApiRoutes(r *gin.Engine) {

	// 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
	var route *gin.RouterGroup

	if len(config.Get[string]("app.api_domain")) == 0 {
		route = r.Group("/api")
	}

	route.Use(middlewares.Common(), middlewares.LimitIP("200-H"))

	authGroup := route.Group("/auth")

	authGroup.Use(middlewares.LimitIP("1000-H"))
	{
		suc := new(clientAuth.SignupController)
		authGroup.POST("/signup/email/exist", suc.IsEmailExist)
		authGroup.POST("/signup/using-email", suc.SignupUsingEmail)

		sic := new(clientAuth.SigninController)
		authGroup.POST("/signin/using-password", sic.SignInByPassword)
		authGroup.POST("/signin/refresh_token", middlewares.AuthJWT(), sic.RefreshToken)

		// 图片验证码
		// authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("50-H"), vcc.ShowCaptcha)

	}

	clientProduct := new(client.ProductController)
	route.GET("/home/products", clientProduct.HomeProducts)

	adminGroup := route.Group("/admin")
	{
		// adminAuthGroup := adminGroup.Group("/auth")
		// {
		// 	suc := new(adminAuth.SignupController)
		// 	adminAuthGroup.POST("/signup/email/exist", suc.SignupUsingEmail)
		// }

		pc := new(admin.ProductController)
		adminGroup.POST("/product/1688", pc.GetProductFrom1688)

	}

}
