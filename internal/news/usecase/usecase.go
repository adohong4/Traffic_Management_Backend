package usecase

import (
	"context"
	"time"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/internal/news"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type newsUC struct {
	cfg      *config.Config
	newsRepo news.Repository
	logger   logger.Logger
}

func NewNewsUseCase(cfg *config.Config, newsRepo news.Repository, log logger.Logger) news.UseCase {
	return &newsUC{cfg: cfg, newsRepo: newsRepo, logger: log}
}

func (uc *newsUC) Create(ctx context.Context, n *models.News) (*models.News, error) {
	if err := n.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}
	n.CreatorId = user.Id

	if err := utils.ValidateStruct(ctx, n); err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	return uc.newsRepo.Create(ctx, n)
}

func (uc *newsUC) Update(ctx context.Context, n *models.News) (*models.News, error) {
	if err := n.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}
	n.ModifierID = &user.Id

	if err := utils.ValidateStruct(ctx, n); err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	return uc.newsRepo.Update(ctx, n)
}

func (uc *newsUC) Delete(ctx context.Context, db *models.News) (*models.News, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "notificationUC.UpdateNotification.GetUserFromCtx"))
	}

	db.ModifierID = &user.Id
	db.UpdatedAt = time.Now()

	if err := utils.ValidateStruct(ctx, db); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "newsUC.UpdateNews.ValidateStruct"))
	}

	result, err := uc.newsRepo.DeleteNews(ctx, db)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (uc *newsUC) FindById(ctx context.Context, id uuid.UUID) (*models.News, error) {
	dt, err := uc.newsRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

func (uc *newsUC) FindAll(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	return uc.newsRepo.FindAll(ctx, pq)
}
