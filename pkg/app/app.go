// Package app 应用信息
package app

import (
	"api/pkg/config"
)

func IsLocal() bool {
	return config.Get[string]("app.env") == "local"
}

func IsProduction() bool {
	return config.Get[string]("app.env") == "production"
}

func IsTesting() bool {
	return config.Get[string]("app.env") == "testing"
}
