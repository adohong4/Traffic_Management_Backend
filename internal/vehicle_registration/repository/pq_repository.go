package repository

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	vehiclelicense "github.com/adohong4/driving-license/internal/vehicle_registration"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Vehicle Document Repository
type vehicleDocRepo struct {
	db *sqlx.DB
}

// Vehicle document new constructor
func NewVehicleDocRepository(db *sqlx.DB) vehiclelicense.Repository {
	return &vehicleDocRepo{db: db}
}

func (r *vehicleDocRepo) CreateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.QueryRowxContext(ctx, createLicenseQuery,
		veDoc.ID, veDoc.OwnerID, veDoc.Brand, veDoc.TypeVehicle, veDoc.VehiclePlateNo, veDoc.ColorPlate, veDoc.ChassisNo, veDoc.EngineNo, veDoc.ColorVehicle,
		veDoc.OwnerName, veDoc.Seats, veDoc.IssueDate, veDoc.ExpiryDate, veDoc.Issuer,
		veDoc.Status, veDoc.Version, veDoc.CreatorId, veDoc.ModifierId, veDoc.CreatedAt, veDoc.UpdatedAt, veDoc.Active,
	).StructScan(v); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.CreateVehicleDoc.StructScan")
	}
	return v, nil
}

func (r *vehicleDocRepo) UpdateVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.QueryRowxContext(ctx, updateLicenseQuery,
		veDoc.OwnerID, veDoc.Brand, veDoc.TypeVehicle, veDoc.VehiclePlateNo, veDoc.ColorPlate, veDoc.ChassisNo, veDoc.EngineNo,
		veDoc.ColorVehicle, veDoc.OwnerName, veDoc.Seats, veDoc.IssueDate, veDoc.ExpiryDate, veDoc.Issuer,
		veDoc.Status, veDoc.ModifierId, veDoc.Active, veDoc.ID,
	).StructScan(v); err != nil {
		return nil, errors.Wrap(err, "vehicleDocRepo.UpdateVehicleDoc.StructScan")
	}
	return v, nil
}

func (r *vehicleDocRepo) DeleteVehicleDoc(ctx context.Context, veDoc *models.VehicleRegistration) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.QueryRowxContext(ctx, deleteLicenseQuery,
		veDoc.Active, veDoc.Version, veDoc.ModifierId, veDoc.UpdatedAt,
	).StructScan(v); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.DeleteVehicle.StructScan")
	}
	return v, nil
}

func (r *vehicleDocRepo) GetVehicleDocs(ctx context.Context, pq *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.VehicleRegistrationList{
			TotalCount:      totalCount,
			TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
			Page:            pq.GetPage(),
			Size:            pq.GetSize(),
			HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			VehicleDocument: make([]*models.VehicleRegistration, 0),
		}, nil
	}

	var NewVehicleDocs = make([]*models.VehicleRegistration, 0, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getVehicleDocuments, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.GetVehicleDocs.QueryRowxContext")
	}
	defer rows.Close()
	for rows.Next() {
		n := &models.VehicleRegistration{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "newsRepo.GetNews.StructScan")
		}
		NewVehicleDocs = append(NewVehicleDocs, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "newsRepo.GetNews.rows.Err")
	}
	return &models.VehicleRegistrationList{
		TotalCount:      totalCount,
		TotalPages:      utils.GetTotalPage(totalCount, pq.GetSize()),
		Page:            pq.GetPage(),
		Size:            pq.GetSize(),
		HasMore:         utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		VehicleDocument: NewVehicleDocs,
	}, nil
}

func (r *vehicleDocRepo) GetVehicleByID(ctx context.Context, vehicleID uuid.UUID) (*models.VehicleRegistration, error) {
	v := &models.VehicleRegistration{}
	if err := r.db.GetContext(ctx, v, getLicenseQuery, vehicleID); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.GetVehicleByID.GetContext")
	}
	return v, nil
}

func (r *vehicleDocRepo) FindByVehiclePlateNO(ctx context.Context, vePlaNO string, query *utils.PaginationQuery) (*models.VehicleRegistrationList, error) {
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findByVehiclePlateNOCount, vePlaNO); err != nil {
		return nil, errors.Wrap(err, "VehicleDocRepo.FindByVehiclePlateNOCount.GetContext")
	}

	if totalCount == 0 {
		return &models.VehicleRegistrationList{
			TotalCount:      totalCount,
			TotalPages:      utils.GetTotalPage(totalCount, query.GetSize()),
			Page:            query.GetPage(),
			Size:            query.GetSize(),
			HasMore:         utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			VehicleDocument: make([]*models.VehicleRegistration, 0),
		}, nil
	}

	var NewVehicleDocs = make([]*models.VehicleRegistration, 0, query.GetSize())
	rows, err := r.db.QueryxContext(ctx, findByVehiclePlateNO, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "NewVehicleDocs.FindByVehiclePlateNOCount.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		n := &models.VehicleRegistration{}
		if err = rows.StructScan(n); err != nil {
			return nil, errors.Wrap(err, "NewVehicleDocs.FindByVehiclePlateNOCount.StructScan")
		}
		NewVehicleDocs = append(NewVehicleDocs, n)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "NewVehicleDocs.FindByVehiclePlateNOCount.rows.err")
	}

	return &models.VehicleRegistrationList{
		TotalCount:      totalCount,
		TotalPages:      utils.GetTotalPage(totalCount, query.GetSize()),
		Page:            query.GetPage(),
		Size:            query.GetSize(),
		HasMore:         utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		VehicleDocument: NewVehicleDocs,
	}, nil
}
