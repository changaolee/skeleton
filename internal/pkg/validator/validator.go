package validator

import (
	"github.com/changaolee/skeleton/pkg/validation"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// validateUsername 检查 username 的合法性.
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if errs := validation.IsQualifiedName(username); len(errs) > 0 {
		return false
	}

	return true
}

// validatePassword 检查 password 的合法性.
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if err := validation.IsValidPassword(password); err != nil {
		return false
	}

	return true
}

// 注册自定义 validation.
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("username", validateUsername)
		_ = v.RegisterValidation("password", validatePassword)
	}
}
