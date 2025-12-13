package news

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(ctx context.Context, n *models.News) (*models.News, error)
	Update(ctx context.Context, n *models.News) (*models.News, error)
	Delete(ctx context.Context, db *models.News) (*models.News, error)
	FindById(ctx context.Context, id uuid.UUID) (*models.News, error)
	FindAll(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error)
}
