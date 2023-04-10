package store

import (
	"context"

	"github.com/changaolee/skeleton/internal/pkg/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.User) error
}
