package v1

import (
	"context"

	mu "github.com/changaolee/skeleton/internal/pkg/model/user"
	"github.com/changaolee/skeleton/internal/pkg/rest"
)

type UsersGetter interface {
	Users() UserInterface
}

type UserInterface interface {
	Create(ctx context.Context, user *mu.User) (*mu.User, error)
	Get(ctx context.Context, name string) (*mu.User, error)
}

type users struct {
	client rest.Interface
}

var _ UserInterface = (*users)(nil)

func newUsers(c *APIV1Client) *users {
	return &users{
		client: c.RESTClient(),
	}
}

func (u *users) Create(ctx context.Context, user *mu.User) (result *mu.User, err error) {
	result = &mu.User{}
	err = u.client.Post().
		AbsPath("/v1/users").
		Body(user).
		Do(ctx).
		Into(result)

	return
}

func (u *users) Get(ctx context.Context, name string) (result *mu.User, err error) {
	result = &mu.User{}
	err = u.client.Get().
		AbsPath("/v1/users/" + name).
		Do(ctx).
		Into(result)

	return
}
