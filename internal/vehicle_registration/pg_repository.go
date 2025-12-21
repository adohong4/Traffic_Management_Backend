package vehicleRegistration

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type Repository interface {
	CreateVehicleDoc(ctx context.Context, vehicleDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
	UpdateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
	ConfirmBlockchainStorage(ctx context.Context, v *models.VehicleRegistration) (*models.VehicleRegistration, error)
	DeleteVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
	GetVehicleDocs(ctx context.Context, pq *utils.PaginationQuery) (*models.VehicleRegistrationList, error)
	GetVehicleByID(ctx context.Context, vehicleID uuid.UUID) (*models.VehicleRegistration, error)
	SearchByVehiclePlateNO(ctx context.Context, vePlaNO string, query *utils.PaginationQuery) (*models.VehicleRegistrationList, error)
	FindVehiclePlateNO(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
	GetCountByType(ctx context.Context) ([]*models.CountItem, error)
	GetTopBrands(ctx context.Context) ([]*models.CountItem, error)
	GetRegistrationStatusStats(ctx context.Context) (*models.StatusCounts, error)
}
