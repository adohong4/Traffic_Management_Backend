package auth

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateUser(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, Id uuid.UUID, modifierId uuid.UUID, version int) error
	GetByID(ctx context.Context, Id uuid.UUID) (*models.User, error)
	FindByIdentity(ctx context.Context, identity string, query *utils.PaginationQuery) (*models.UsersList, error)
	GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	ConnectWallet(ctx context.Context, user *models.User) (*models.UserWithToken, error)
}
