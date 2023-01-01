package user

import (
	"github.com/changaolee/skeleton/internal/skeleton/biz"
	"github.com/changaolee/skeleton/internal/skeleton/store"
)

// UserController 是 user 模块在 Controller 层的实现，用来处理用户模块的请求
type UserController struct {
	b biz.IBiz
}

// New 创建一个 user controller.
func New(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}
