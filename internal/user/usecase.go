package user

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
)

type UseCase interface {
	CreateUser(ctx context.Context, db *models.User) (*models.User, error)
}
