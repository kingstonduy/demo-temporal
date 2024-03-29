package data

import (
	"context"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/database"
)

type UserRepository interface {
	FindByToken(ctx context.Context, token string) (*domain.User, error)
	CountById(ctx context.Context, id string) (int64, error)
	CreateUser(ctx context.Context, user *domain.User) error
	FindUserById(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	database.EnableTransactor
}
