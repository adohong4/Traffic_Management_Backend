package notification

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error)
	UpdateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error)
	DeleteNotification(ctx context.Context, db *models.Notification) (*models.Notification, error)
	GetNotification(ctx context.Context, pq *utils.PaginationQuery) (*models.NotificationList, error)
	GetNotificationByID(ctx context.Context, Id uuid.UUID) (*models.Notification, error)
	SearchNotificationByTitle(ctx context.Context, title string, pq *utils.PaginationQuery) (*models.NotificationList, error)
}
