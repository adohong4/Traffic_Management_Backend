package repository

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	notification "github.com/adohong4/driving-license/internal/notification"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type notificationRepo struct {
	db *sqlx.DB
}

func NewNotificationRepo(db *sqlx.DB) notification.Repository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) CreateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error) {
	n := &models.Notification{}
	if err := r.db.QueryRowxContext(ctx, createNotificationQuery,
		db.Id, db.Title, db.Content, db.Target, db.CreatorId, db.CreatedAt, db.UpdatedAt, db.Active,
	).StructScan(n); err != nil {
		return nil, errors.Wrap(err, "notificationRepo.CreateNotification.StructScan")
	}
	return n, nil
}

func (r *notificationRepo) UpdateNotification(ctx context.Context, db *models.Notification) (*models.Notification, error) {
	n := &models.Notification{}
	if err := r.db.QueryRowxContext(ctx, updateNotificationQuery,
		db.Title, db.Content, db.Target, db.ModifierID, db.UpdatedAt, db.Active, db.Id,
	).StructScan(n); err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.UpdateNotification.StructScan")
	}
	return n, nil
}

func (r *notificationRepo) DeleteNotification(ctx context.Context, db *models.Notification) (*models.Notification, error) {
	n := &models.Notification{}
	if err := r.db.QueryRowxContext(ctx, deleteNotificationQuery,
		db.ModifierID, db.UpdatedAt, db.Id,
	).StructScan(n); err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.DeleteNotification.StructScan")
	}
	return n, nil
}

func (r *notificationRepo) GetNotification(ctx context.Context, pq *utils.PaginationQuery) (*models.NotificationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.GetNotification.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.NotificationList{
			TotalCount:   totalCount,
			TotalPages:   utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:         pq.GetPage(),
			Size:         pq.GetSize(),
			HasMore:      utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Notification: make([]*models.Notification, 0),
		}, nil
	}

	var newNotifications = make([]*models.Notification, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getNotification, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.GetNotification.QueryxContext")
	}
	defer rows.Close()
	for rows.Next() {
		n := &models.Notification{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "NotificationRepo.GetNotification.StructScan")
		}
		newNotifications = append(newNotifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.GetNotification.rows.Err")
	}

	return &models.NotificationList{
		TotalCount:   totalCount,
		TotalPages:   utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:         pq.GetPage(),
		Size:         pq.GetSize(),
		HasMore:      utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Notification: newNotifications,
	}, nil
}

func (r *notificationRepo) GetNotificationByID(ctx context.Context, notificationID string) (*models.Notification, error) {
	n := &models.Notification{}
	if err := r.db.GetContext(ctx, n, getNotificationByIdQuery, notificationID); err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.GetNotificationByID.GetContext")
	}
	return n, nil
}

func (r *notificationRepo) SearchNotificationByTitle(ctx context.Context, title string, pq *utils.PaginationQuery) (*models.NotificationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findByTitleCount, title); err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.SearchNotificationByTitle.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.NotificationList{
			TotalCount:   totalCount,
			TotalPages:   utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:         pq.GetPage(),
			Size:         pq.GetSize(),
			HasMore:      utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Notification: make([]*models.Notification, 0),
		}, nil
	}

	var NewNotifications = make([]*models.Notification, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, searchByTitleQuery, title, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.SearchNotificationByTitle.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.Notification{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "NotificationRepo.SearchNotificationByTitle.StructScan")
		}
		NewNotifications = append(NewNotifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NotificationRepo.SearchNotificationByTitle.rows.Err")
	}

	return &models.NotificationList{
		TotalCount:   totalCount,
		TotalPages:   utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:         pq.GetPage(),
		Size:         pq.GetSize(),
		HasMore:      utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Notification: NewNotifications,
	}, nil
}
