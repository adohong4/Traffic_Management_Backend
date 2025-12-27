package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/internal/notification"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type notificationUC struct {
	cfg              *config.Config
	notificationRepo notification.Repository
	logger           logger.Logger
}

func NewNotificationUseCase(cfg *config.Config, notificationRepo notification.Repository, log logger.Logger) notification.UseCase {
	return &notificationUC{cfg: cfg, notificationRepo: notificationRepo, logger: log}
}

func (n *notificationUC) CreateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error) {
	if err := db.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "notificationUC.CreateNotification.PrepareCreate"))
	}

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "notificationUC.CreateNotification.GetUserFromCtx"))
	}

	db.CreatorId = user.Id

	if err = utils.ValidateStruct(ctx, db); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "notificationUC.CreateNotification.ValidateStruct"))
	}

	notificationResult, err := n.notificationRepo.CreateNotification(ctx, db)
	if err != nil {
		return nil, err
	}

	return notificationResult, nil
}

func (n *notificationUC) UpdateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "notificationUC.UpdateNotification.GetUserFromCtx"))
	}

	db.ModifierID = &user.Id
	db.UpdatedAt = time.Now()

	if err := utils.ValidateStruct(ctx, db); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "notificationUC.UpdateNotification.ValidateStruct"))
	}

	notificationResult, err := n.notificationRepo.UpdateNotification(ctx, db)
	if err != nil {
		return nil, err
	}
	return notificationResult, nil
}

func (n *notificationUC) DeleteNotification(ctx context.Context, db *models.Notification) (*models.Notification, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "notificationUC.UpdateNotification.GetUserFromCtx"))
	}

	db.ModifierID = &user.Id
	db.UpdatedAt = time.Now()

	if err := utils.ValidateStruct(ctx, db); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "notificationUC.UpdateNotification.ValidateStruct"))
	}

	notificationResult, err := n.notificationRepo.DeleteNotification(ctx, db)
	if err != nil {
		return nil, err
	}
	return notificationResult, nil
}

func (n *notificationUC) GetNotification(ctx context.Context, pq *utils.PaginationQuery) (*models.NotificationList, error) {
	return n.notificationRepo.GetNotification(ctx, pq)
}

func (n *notificationUC) GetNotificationByID(ctx context.Context, Id uuid.UUID) (*models.Notification, error) {
	dt, err := n.notificationRepo.GetNotificationByID(ctx, Id)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

func (n *notificationUC) SearchNotificationByTitle(ctx context.Context, title string, pq *utils.PaginationQuery) (*models.NotificationList, error) {
	return n.notificationRepo.SearchNotificationByTitle(ctx, title, pq)
}

func (n *notificationUC) GetMyNotifications(ctx context.Context, pq *utils.PaginationQuery) (*models.NotificationList, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}

	// Giả sử user có IdentityNo (CCCD)
	if user.IdentityNo == "" {
		return nil, httpErrors.NewBadRequestError(errors.New("user identity number is missing"))
	}

	return n.notificationRepo.GetNotificationsForUser(ctx, user.CreatedAt, user.IdentityNo, pq)
}

func (n *notificationUC) GetMyNotificationByID(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}

	if user.IdentityNo == "" {
		return nil, httpErrors.NewBadRequestError(errors.New("user identity number is missing"))
	}

	// Nếu là personal → đánh dấu đã đọc và trả về
	updatedNoti, err := n.notificationRepo.MarkAsReadAndGet(ctx, notificationID, user.IdentityNo)
	if err != nil {
		return nil, err
	}
	if updatedNoti != nil {
		return updatedNoti, nil
	}

	// Nếu không phải personal → chỉ lấy thông thường (nếu thuộc user hoặc all)
	// Ta lấy thông báo nếu nó là 'all' hoặc 'personal' thuộc user
	noti, err := n.notificationRepo.GetNotificationByID(ctx, notificationID)
	if err != nil {
		return nil, err
	}

	// Kiểm tra quyền xem
	if noti.Target == "personal" && noti.TargetUser != user.IdentityNo {
		return nil, httpErrors.NewRestError(http.StatusForbidden, "you are not allowed to view this notification", nil)
	}

	return noti, nil
}
