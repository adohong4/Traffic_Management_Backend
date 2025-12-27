package repository

import (
	"context"
	"database/sql"

	"github.com/adohong4/driving-license/internal/models"
	trafficviolation "github.com/adohong4/driving-license/internal/traffic_violation"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type TrafficViolationRepo struct {
	db *sqlx.DB
}

func NewTrafficViolationRepo(db *sqlx.DB) trafficviolation.Repository {
	return &TrafficViolationRepo{db: db}
}

func (r *TrafficViolationRepo) CreateTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error) {
	t := &models.TrafficViolation{}
	if err := r.db.QueryRowxContext(ctx, createTrafficViolationQuery,
		tv.Id, tv.VehiclePlateNo, tv.Date, tv.Type, tv.Address, tv.Description, tv.Points, tv.FineAmount, tv.ExpiryDate,
		tv.Status, tv.Version, tv.CreatorId, tv.ModifierId, tv.CreatedAt, tv.UpdatedAt, tv.Active,
	).StructScan(t); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.CreateTrafficViolation.StructScan")
	}
	return t, nil
}

func (r *TrafficViolationRepo) UpdateTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error) {
	t := &models.TrafficViolation{}
	if err := r.db.QueryRowxContext(ctx, updateTrafficViolationQuery,
		tv.VehiclePlateNo, tv.Date, tv.Type, tv.Address, tv.Description, tv.Points, tv.FineAmount, tv.ExpiryDate,
		tv.Status, tv.ModifierId, tv.UpdatedAt, tv.Id,
	).StructScan(t); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.UpdateTrafficViolation.StructScan")
	}
	return t, nil
}

func (r *TrafficViolationRepo) DeleteTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error) {
	t := &models.TrafficViolation{}
	if err := r.db.QueryRowxContext(ctx, deleteTrafficViolationQuery, tv.ModifierId, tv.UpdatedAt, tv.Id).StructScan(t); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.DeleteTrafficViolation.StructScan")
	}
	return t, nil
}

func (r *TrafficViolationRepo) GetTrafficViolationById(ctx context.Context, Id uuid.UUID) (*models.TrafficViolation, error) {
	t := &models.TrafficViolation{}
	if err := r.db.GetContext(ctx, t, getTrafficViolationByIdQuery, Id); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.DeleteTrafficViolation.StructScan")
	}
	return t, nil
}

func (r *TrafficViolationRepo) GetAllTrafficViolation(ctx context.Context, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTrafficViolationTotalCount); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.GetAllTrafficViolation.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.TrafficViolationList{
			TotalCount:       totalCount,
			TotalPages:       utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:             pq.GetPage(),
			Size:             pq.GetSize(),
			HasMore:          utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			TrafficViolation: make([]*models.TrafficViolation, 0),
		}, nil
	}

	var NewTrafficViolation = make([]*models.TrafficViolation, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getTrafficViolationQuery, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.GetAllTrafficViolation.NewTrafficViolation")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.TrafficViolation{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "TrafficViolationRepo.GetAllTrafficViolation.StructScan")
		}
		NewTrafficViolation = append(NewTrafficViolation, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.rows.err")
	}

	return &models.TrafficViolationList{
		TotalCount:       totalCount,
		TotalPages:       utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:             pq.GetPage(),
		Size:             pq.GetSize(),
		HasMore:          utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		TrafficViolation: NewTrafficViolation,
	}, nil
}

func (r *TrafficViolationRepo) SearchTrafficViolation(ctx context.Context, vpn string, query *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findVehiclePlateNoCount, vpn); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.GetAllTrafficViolation.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.TrafficViolationList{
			TotalCount:       totalCount,
			TotalPages:       utils.GetTotalPage(totalCount, query.GetSize()),
			Page:             query.GetPage(),
			Size:             query.GetSize(),
			HasMore:          utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			TrafficViolation: make([]*models.TrafficViolation, 0),
		}, nil
	}

	var NewTrafficViolation = make([]*models.TrafficViolation, 0, query.GetSize())
	rows, err := r.db.QueryxContext(ctx, searchByVehicleNo, vpn, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.GetAllTrafficViolation.NewTrafficViolation")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.TrafficViolation{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "TrafficViolationRepo.GetAllTrafficViolation.StructScan")
		}
		NewTrafficViolation = append(NewTrafficViolation, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.rows.err")
	}

	return &models.TrafficViolationList{
		TotalCount:       totalCount,
		TotalPages:       utils.GetTotalPage(totalCount, query.GetSize()),
		Page:             query.GetPage(),
		Size:             query.GetSize(),
		HasMore:          utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		TrafficViolation: NewTrafficViolation,
	}, nil
}

func (r *TrafficViolationRepo) GetTrafficViolationStats(ctx context.Context) (*models.TrafficViolationStats, error) {
	var stats models.TrafficViolationStats
	err := r.db.GetContext(ctx, &stats, getTrafficViolationStatsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.GetTrafficViolationStats")
	}
	return &stats, nil
}

func (r *TrafficViolationRepo) GetTrafficViolationStatusStats(ctx context.Context) ([]*models.TrafficViolationStatusStats, error) {
	var stats []*models.TrafficViolationStatusStats

	err := r.db.SelectContext(ctx, &stats, getTrafficViolationStatusStatsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.GetTrafficViolationStatusStats.SelectContext")
	}

	return stats, nil
}

func (r *TrafficViolationRepo) GetViolationsByVehiclePlateNo(ctx context.Context, plateNo string, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, getTotalByPlateNo, plateNo); err != nil {
		return nil, errors.Wrap(err, "GetViolationsByVehiclePlateNo.total")
	}

	list := &models.TrafficViolationList{
		TotalCount:       total,
		TotalPages:       utils.GetTotalPage(total, pq.GetSize()),
		Page:             pq.GetPage(),
		Size:             pq.GetSize(),
		HasMore:          utils.GetHasMore(pq.GetPage(), total, pq.GetSize()),
		TrafficViolation: []*models.TrafficViolation{},
	}

	if total == 0 {
		return list, nil
	}

	var items []*models.TrafficViolation
	if err := r.db.SelectContext(ctx, &items, getViolationsByPlateNo, plateNo, pq.GetOffset(), pq.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "GetViolationsByVehiclePlateNo.Select")
	}

	list.TrafficViolation = items
	return list, nil
}

func (r *TrafficViolationRepo) GetMyViolationsByOwnerID(ctx context.Context, ownerID uuid.UUID, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, getTotalViolationsByOwnerID, ownerID); err != nil {
		return nil, errors.Wrap(err, "GetMyViolationsByOwnerID.total")
	}

	list := &models.TrafficViolationList{
		TotalCount:       total,
		TotalPages:       utils.GetTotalPage(total, pq.GetSize()),
		Page:             pq.GetPage(),
		Size:             pq.GetSize(),
		HasMore:          utils.GetHasMore(pq.GetPage(), total, pq.GetSize()),
		TrafficViolation: []*models.TrafficViolation{},
	}

	if total == 0 {
		return list, nil
	}

	var items []*models.TrafficViolation
	if err := r.db.SelectContext(ctx, &items, getViolationsByOwnerID, ownerID, pq.GetOffset(), pq.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "GetMyViolationsByOwnerID.Select")
	}

	list.TrafficViolation = items
	return list, nil
}

func (r *TrafficViolationRepo) GetMyViolationsByWallet(ctx context.Context, wallet string, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, getTotalViolationsByWallet, wallet); err != nil {
		return nil, errors.Wrap(err, "GetMyViolationsByWallet.total")
	}

	list := &models.TrafficViolationList{
		TotalCount:       total,
		TotalPages:       utils.GetTotalPage(total, pq.GetSize()),
		Page:             pq.GetPage(),
		Size:             pq.GetSize(),
		HasMore:          utils.GetHasMore(pq.GetPage(), total, pq.GetSize()),
		TrafficViolation: []*models.TrafficViolation{},
	}

	if total == 0 {
		return list, nil
	}

	var items []*models.TrafficViolation
	if err := r.db.SelectContext(ctx, &items, getViolationsByWallet, wallet, pq.GetOffset(), pq.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "GetMyViolationsByWallet.Select")
	}

	list.TrafficViolation = items
	return list, nil
}

func (r *TrafficViolationRepo) GetVehiclePlateNoIfOwned(ctx context.Context, vehicleID, ownerID uuid.UUID) (string, error) {
	var plateNo string
	err := r.db.GetContext(ctx, &plateNo,
		`SELECT vehicle_no FROM vehicle_registration WHERE id = $1 AND owner_id = $2 AND active = true`,
		vehicleID, ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", sql.ErrNoRows
		}
		return "", errors.Wrap(err, "TrafficViolationRepo.GetVehiclePlateNoIfOwned")
	}
	return plateNo, nil
}

func (r *TrafficViolationRepo) GetTrafficViolationByIDAndOwnerID(ctx context.Context, violationID, ownerID uuid.UUID) (*models.TrafficViolation, error) {
	v := &models.TrafficViolation{}
	err := r.db.GetContext(ctx, v, getTrafficViolationByIDAndOwner, violationID, ownerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.GetTrafficViolationByIDAndOwnerID.GetContext")
	}
	return v, nil
}

func (r *TrafficViolationRepo) GetViolationsByLicenseWallet(ctx context.Context, wallet string, pq *utils.PaginationQuery) (*models.TrafficViolationList, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, getTotalViolationsByLicenseWallet, wallet); err != nil {
		return nil, errors.Wrap(err, "GetViolationsByLicenseWallet.total")
	}

	list := &models.TrafficViolationList{
		TotalCount:       total,
		TotalPages:       utils.GetTotalPage(total, pq.GetSize()),
		Page:             pq.GetPage(),
		Size:             pq.GetSize(),
		HasMore:          utils.GetHasMore(pq.GetPage(), total, pq.GetSize()),
		TrafficViolation: []*models.TrafficViolation{},
	}

	if total == 0 {
		return list, nil
	}

	var items []*models.TrafficViolation
	if err := r.db.SelectContext(ctx, &items, getViolationsByLicenseWallet, wallet, pq.GetOffset(), pq.GetLimit()); err != nil {
		return nil, errors.Wrap(err, "GetViolationsByLicenseWallet.Select")
	}

	list.TrafficViolation = items
	return list, nil
}
