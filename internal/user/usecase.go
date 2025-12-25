package user

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
)

type UseCase interface {
	CreateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error)
}
