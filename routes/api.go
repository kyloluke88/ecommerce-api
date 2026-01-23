package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterApiRoutes 注册网页相关路由
func RegisterApiRoutes(r *gin.Engine) {

	fmt.Println("注册 API 路由")
	// 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
	v1 := r.Group("/v1")
	// 注册一个路由
	v1.GET("/", func(c *gin.Context) {
		// 以 JSON 格式响应
		c.JSON(http.StatusOK, gin.H{
			"Hello": "World!",
		})
	})

	v1.GET("/health", func(c *gin.Context) {

		c.Writer.Write([]byte("ok"))
	})
}
