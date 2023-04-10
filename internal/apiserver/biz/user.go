package biz

import (
	"context"
	"regexp"

	"github.com/changaolee/skeleton/internal/apiserver/store"
	"github.com/changaolee/skeleton/internal/pkg/errno"
	"github.com/changaolee/skeleton/internal/pkg/model"
)

type UserBiz interface {
	Create(ctx context.Context, user *model.User) error
}

type userBiz struct {
	s store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func newUsers(b *biz) *userBiz {
	return &userBiz{s: b.s}
}

func (b *userBiz) Create(ctx context.Context, user *model.User) error {
	if err := b.s.Users().Create(ctx, user); err != nil {
		if matched, _ := regexp.MatchString("", err.Error()); matched {
			return errno.ErrUserAlreadyExist
		}
		return err
	}
	return nil
}
