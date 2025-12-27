package trafficviolation

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type Repository interface {
	CreateTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error)
	UpdateTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error)
	DeleteTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error)
	GetTrafficViolationById(ctx context.Context, Id uuid.UUID) (*models.TrafficViolation, error)
	GetAllTrafficViolation(ctx context.Context, pq *utils.PaginationQuery) (*models.TrafficViolationList, error)
	SearchTrafficViolation(ctx context.Context, vpn string, query *utils.PaginationQuery) (*models.TrafficViolationList, error)
	GetTrafficViolationStats(ctx context.Context) (*models.TrafficViolationStats, error)
	GetTrafficViolationStatusStats(ctx context.Context) ([]*models.TrafficViolationStatusStats, error)
	GetViolationsByVehiclePlateNo(ctx context.Context, plateNo string, pq *utils.PaginationQuery) (*models.TrafficViolationList, error)
	GetMyViolationsByOwnerID(ctx context.Context, ownerID uuid.UUID, pq *utils.PaginationQuery) (*models.TrafficViolationList, error)
	GetMyViolationsByWallet(ctx context.Context, wallet string, pq *utils.PaginationQuery) (*models.TrafficViolationList, error)
	GetVehiclePlateNoIfOwned(ctx context.Context, vehicleID, ownerID uuid.UUID) (string, error)
	GetTrafficViolationByIDAndOwnerID(ctx context.Context, violationID, ownerID uuid.UUID) (*models.TrafficViolation, error)
	GetViolationsByLicenseWallet(ctx context.Context, wallet string, pq *utils.PaginationQuery) (*models.TrafficViolationList, error)
}
