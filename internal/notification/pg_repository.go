package notification

import (
	"context"
	"time"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type Repository interface {
	CreateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error)
	UpdateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error)
	DeleteNotification(ctx context.Context, db *models.Notification) (*models.Notification, error)
	GetNotification(ctx context.Context, pq *utils.PaginationQuery) (*models.NotificationList, error)
	GetNotificationByID(ctx context.Context, Id uuid.UUID) (*models.Notification, error)
	SearchNotificationByTitle(ctx context.Context, title string, pq *utils.PaginationQuery) (*models.NotificationList, error)
	GetNotificationsForUser(ctx context.Context, userCreatedAt time.Time, identityNo string, pq *utils.PaginationQuery) (*models.NotificationList, error)
	MarkAsReadAndGet(ctx context.Context, notificationID uuid.UUID, identityNo string) (*models.Notification, error)
}
