package repository

import (
	"context"

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
		tv.Id, tv.VehiclePlateNo, tv.Date, tv.Type, tv.Description, tv.Points, tv.FineAmount,
		tv.Status, tv.Version, tv.CreatorId, tv.ModifierId, tv.CreatedAt, tv.UpdatedAt, tv.Active,
	).StructScan(t); err != nil {
		return nil, errors.Wrap(err, "TrafficViolationRepo.CreateTrafficViolation.StructScan")
	}
	return t, nil
}

func (r *TrafficViolationRepo) UpdateTrafficViolation(ctx context.Context, tv *models.TrafficViolation) (*models.TrafficViolation, error) {
	t := &models.TrafficViolation{}
	if err := r.db.QueryRowxContext(ctx, updateTrafficViolationQuery,
		tv.VehiclePlateNo, tv.Date, tv.Type, tv.Description, tv.Points, tv.FineAmount,
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
