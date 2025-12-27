package usecase

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	trafficviolation "github.com/adohong4/driving-license/internal/traffic_violation"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type TrafficViolationUC struct {
	cfg                  *config.Config
	TrafficViolationRepo trafficviolation.Repository
	logger               logger.Logger
}

func NewTrafficViolationUseCase(cfg *config.Config, TrafficViolationRepo trafficviolation.Repository, logger logger.Logger) trafficviolation.UseCase {
	return &TrafficViolationUC{cfg: cfg, TrafficViolationRepo: TrafficViolationRepo, logger: logger}
}

func (u *TrafficViolationUC) CreateTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error) {
	if err := tv.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "TrafficViolationUC.CreateTrafficViolation.PrepareCreate"))
	}

	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "TrafficViolationUC.Create.GetUserFromCtx"))
	}

	tv.CreatorId = user.Id

	if err = utils.ValidateStruct(ctx, tv); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "TrafficViolationUC.Create.ValidateStruct"))
	}

	n, err := u.TrafficViolationRepo.CreateTrafficViolation(ctx, tv)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (u *TrafficViolationUC) UpdateTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "TrafficViolationUC.UpdateTrafficViolation.GetUserFromCtx"))
	}

	tv.ModifierId = &user.Id

	if err = tv.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "TrafficViolationUC.UpdateTrafficViolation.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, tv); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "TrafficViolationUC.UpdateTrafficViolation.ValidateStruct"))
	}

	updatedLicense, err := u.TrafficViolationRepo.UpdateTrafficViolation(ctx, tv)
	if err != nil {
		return nil, err
	}

	return updatedLicense, nil
}

func (u *TrafficViolationUC) DeleteTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "TrafficViolationUC.DeleteTrafficViolation.GetUserFromCtx"))
	}

	tv.ModifierId = &user.Id

	if err = tv.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "TrafficViolationUC.DeleteTrafficViolation.PrepareCreate"))
	}

	if err := utils.ValidateStruct(ctx, tv); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "TrafficViolationUC.DeleteTrafficViolation.ValidateStruct"))
	}

	DeleteReport, err := u.TrafficViolationRepo.DeleteTrafficViolation(ctx, tv)
	if err != nil {
		return nil, err
	}

	return DeleteReport, nil
}

func (u *TrafficViolationUC) GetTrafficViolationById(ctx context.Context, Id uuid.UUID) (*models.TrafficViolation, error) {
	n, err := u.TrafficViolationRepo.GetTrafficViolationById(ctx, Id)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (u *TrafficViolationUC) GetAllTrafficViolation(ctx context.Context, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	return u.TrafficViolationRepo.GetAllTrafficViolation(ctx, pq)
}

func (u *TrafficViolationUC) SearchTrafficViolation(ctx context.Context, vpn string, query *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	return u.TrafficViolationRepo.SearchTrafficViolation(ctx, vpn, query)
}

func (u *TrafficViolationUC) GetTrafficViolationStats(ctx context.Context) (*models.TrafficViolationStats, error) {
	return u.TrafficViolationRepo.GetTrafficViolationStats(ctx)
}

func (u *TrafficViolationUC) GetTrafficViolationStatusStats(ctx context.Context) ([]*models.TrafficViolationStatusStats, error) {
	return u.TrafficViolationRepo.GetTrafficViolationStatusStats(ctx)
}

func (u *TrafficViolationUC) GetMyViolations(ctx context.Context, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}

	if *user.UserAddress != "" {
		return u.TrafficViolationRepo.GetMyViolationsByWallet(ctx, *user.UserAddress, pq)
	}

	return u.TrafficViolationRepo.GetMyViolationsByOwnerID(ctx, user.Id, pq)
}

func (u *TrafficViolationUC) GetViolationsByMyVehicle(ctx context.Context, vehicleID uuid.UUID, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}

	plateNo, err := u.TrafficViolationRepo.GetVehiclePlateNoIfOwned(ctx, vehicleID, user.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, httpErrors.NewRestError(http.StatusNotFound, "vehicle not found or not owned by you", nil)
		}
		return nil, err
	}

	return u.TrafficViolationRepo.GetViolationsByVehiclePlateNo(ctx, plateNo, pq)
}

func (u *TrafficViolationUC) GetMyTrafficViolationByID(ctx context.Context, violationID uuid.UUID) (*models.TrafficViolation, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}

	violation, err := u.TrafficViolationRepo.GetTrafficViolationByIDAndOwnerID(ctx, violationID, user.Id)
	if err != nil {
		return nil, err
	}
	if violation == nil {
		return nil, httpErrors.NewRestError(http.StatusNotFound, "violation not found or not related to your vehicle", nil)
	}
	return violation, nil
}

func (u *TrafficViolationUC) GetViolationsByMyLicense(ctx context.Context, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, httpErrors.NewUnauthorizedError(err)
	}

	// Giả sử user có wallet_address trong context (từ JWT)
	if *user.UserAddress == "" {
		return nil, httpErrors.NewBadRequestError(errors.New("user has no linked wallet address"))
	}

	return u.TrafficViolationRepo.GetViolationsByLicenseWallet(ctx, *user.UserAddress, pq)
}
