package bootstrap

import (
	authRequest "api/app/requests/client/auth"

	"github.com/go-playground/validator/v10"

	"api/app/requests"
)

// Validate 是全局 validator 实例
// 只初始化一次，并发安全
var validate *validator.Validate

// 负责创建 validator 并显式注册各 request 模块的自定义校验
func SetupCustomRules() {
	validate = validator.New()

	// 显式注册 admin 模块的校验规则
	// 单向依赖：bootstrap -> requests/admin
	authRequest.RegisterEmailIsExist(validate)

	// 注入到 requests 包
	requests.Validate = validate
}
