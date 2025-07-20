package repository

import (
	"context"
	"database/sql"

	driverlicense "github.com/adohong4/driving-license/internal/driver_license"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DriverLicenseRepo struct {
	db *sqlx.DB
}

func NewDriverLicenseRepo(db *sqlx.DB) driverlicense.Repository {
	return &DriverLicenseRepo{db: db}
}

func (r *DriverLicenseRepo) CreateDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.QueryRowxContext(ctx, createDriverLicenseQuery,
		dl.Id, dl.Name, dl.DOB, dl.IdentityNo, dl.OwnerAddress, dl.LicenseNo,
		dl.IssueDate, dl.ExpiryDate, dl.Status, dl.LicenseType, dl.AuthorityId, dl.IssuingAuthority,
		dl.Nationality, dl.Point, dl.Version, dl.CreatorId, dl.ModifierId, dl.CreatedAt, dl.UpdatedAt, dl.Active,
	).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.CreateDriverLicense.StructScan")
	}
	return d, nil
}

func (r *DriverLicenseRepo) UpdateDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.QueryRowxContext(ctx, updateDriverLicenseQuery,
		dl.Name, dl.DOB, dl.OwnerAddress, dl.ExpiryDate, dl.Status,
		dl.Nationality, dl.Point, dl.ModifierId, dl.UpdatedAt, dl.Id,
	).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.UpdateDriverLicense.StructScan")
	}
	return d, nil
}

func (r *DriverLicenseRepo) DeleteDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.QueryRowxContext(ctx, deleteDriverLicenseQuery, dl.ModifierId, dl.UpdatedAt, dl.Id).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.DeleteDriverLicense.StructScan")
	}
	return d, nil
}

func (r *DriverLicenseRepo) GetDriverLicense(ctx context.Context, pq *utils.PaginationQuery) (*models.DrivingLicenseList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetDriverLicense.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.DrivingLicenseList{
			TotalCount:     totalCount,
			TotalPages:     utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:           pq.GetPage(),
			Size:           pq.GetSize(),
			HasMore:        utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			DrivingLicense: make([]*models.DrivingLicense, 0),
		}, nil
	}

	var NewDriverLicense = make([]*models.DrivingLicense, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getDriverLicense, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetDriverLicense.NewDriverLicense")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.DrivingLicense{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "DriverLicenseRepo.GetDriverLicense.StructScan")
		}
		NewDriverLicense = append(NewDriverLicense, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.rows.err")
	}

	return &models.DrivingLicenseList{
		TotalCount:     totalCount,
		TotalPages:     utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:           pq.GetPage(),
		Size:           pq.GetSize(),
		HasMore:        utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		DrivingLicense: NewDriverLicense,
	}, nil
}

func (r *DriverLicenseRepo) GetDriverLicenseById(ctx context.Context, Id uuid.UUID) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.GetContext(ctx, d, getDriverLicenseByIdQuery, Id); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetDriverLicenseById.GetContext")
	}
	return d, nil
}

func (r *DriverLicenseRepo) SearchByLicenseNo(ctx context.Context, lno string, query *utils.PaginationQuery) (*models.DrivingLicenseList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findLicenseNOCount, lno); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.SearchByLicenseNo.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.DrivingLicenseList{
			TotalCount:     totalCount,
			TotalPages:     utils.GetTotalPage(totalCount, query.GetSize()),
			Page:           query.GetPage(),
			Size:           query.GetSize(),
			HasMore:        utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			DrivingLicense: make([]*models.DrivingLicense, 0),
		}, nil
	}

	var NewDriverLicense = make([]*models.DrivingLicense, 0, query.GetSize())
	rows, err := r.db.QueryxContext(ctx, searchByLicenseNo, lno, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.SearchByLicenseNo.NewDriverLicense")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.DrivingLicense{}
		if err := rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "DriverLicenseRepo.SearchByLicenseNo.StructScan")
		}
		NewDriverLicense = append(NewDriverLicense, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.rows.err")
	}

	return &models.DrivingLicenseList{
		TotalCount:     totalCount,
		TotalPages:     utils.GetTotalPage(totalCount, query.GetSize()),
		Page:           query.GetPage(),
		Size:           query.GetSize(),
		HasMore:        utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		DrivingLicense: NewDriverLicense,
	}, nil
}

func (r *DriverLicenseRepo) FindLicenseNO(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	err := r.db.QueryRowxContext(ctx, findLicenseNO, dl.LicenseNo).StructScan(d)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.findVehiclePlateNO.QueryRowxContext")
	}
	return d, nil
}
