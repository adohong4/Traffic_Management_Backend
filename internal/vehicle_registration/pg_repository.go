package vehiclelicense

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
)

type Repository interface {
	CreateVehicleDoc(ctx context.Context, vehicleDoc *models.VehicleRegistration) (*models.VehicleRegistration, error)
}
