package driverlicense

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
)

type Repository interface {
	CreateDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error)
	UpdateDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error)
	ConfirmBlockchainStorage(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error)
	UpdateWalletAddress(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error)
	DeleteDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error)
	GetDriverLicense(ctx context.Context, pq *utils.PaginationQuery) (*models.DrivingLicenseList, error)
	GetDriverLicenseById(ctx context.Context, Id uuid.UUID) (*models.DrivingLicense, error)
	GetDriverLicenseByWalletAddress(ctx context.Context, address string) (*models.DrivingLicense, error)
	GetDriverLicenseByLicenseNO(ctx context.Context, address string) (*models.DrivingLicense, error)
	SearchByLicenseNo(ctx context.Context, lno string, query *utils.PaginationQuery) (*models.DrivingLicenseList, error)
	FindLicenseNO(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error)
	GetStatusDistribution(ctx context.Context) (*models.StatusDistributionResponse, error)
	GetLicenseTypeDistribution(ctx context.Context) (*models.LicenseTypeDistributionResponse, error)
	GetLicenseTypeStatusDistribution(ctx context.Context) (*models.LicenseTypeDetailDistributionResponse, error)
	GetCityStatusDistribution(ctx context.Context) (*models.CityDetailDistributionResponse, error)
}
