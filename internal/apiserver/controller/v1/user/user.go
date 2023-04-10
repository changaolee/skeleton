package user

import (
	"github.com/changaolee/skeleton/internal/apiserver/biz"
	"github.com/changaolee/skeleton/internal/apiserver/store"
)

type UserController struct {
	b biz.IBiz
}

// New 创建一个 user controller.
func New(s store.IStore) *UserController {
	return &UserController{b: biz.New(s)}
}
