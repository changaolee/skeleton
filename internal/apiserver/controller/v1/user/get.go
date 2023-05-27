package user

import (
	"github.com/changaolee/skeleton/internal/pkg/core"
	"github.com/changaolee/skeleton/pkg/log"
	"github.com/gin-gonic/gin"
)

// Get 通过用户名查询用户信息.
func (u *UserController) Get(c *gin.Context) {
	log.C(c).Infow("Get user function called.")

	user, err := u.b.Users().Get(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
