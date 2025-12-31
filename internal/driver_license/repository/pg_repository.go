package repository

import (
	"context"
	"database/sql"
	"sort"

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
		dl.Id, dl.Name, dl.Avatar, dl.DOB, dl.IdentityNo, dl.OwnerAddress, dl.OwnerCity, dl.LicenseNo,
		dl.IssueDate, dl.ExpiryDate, dl.Status, dl.LicenseType, dl.AuthorityId, dl.IssuingAuthority,
		dl.Nationality, dl.Point, dl.WalletAddress, dl.OnBlockchain, dl.BlockchainTxHash,
		dl.Version, dl.CreatorId, dl.ModifierId, dl.CreatedAt, dl.UpdatedAt, dl.Active,
	).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.CreateDriverLicense.StructScan")
	}
	return d, nil
}

func (r *DriverLicenseRepo) UpdateDriverLicense(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.QueryRowxContext(ctx, updateDriverLicenseQuery,
		dl.Name, dl.DOB, dl.OwnerAddress, dl.OwnerCity, dl.ExpiryDate, dl.Status,
		dl.Nationality, dl.Point, dl.ModifierId, dl.UpdatedAt, dl.Id,
	).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.UpdateDriverLicense.StructScan")
	}
	return d, nil
}

func (r *DriverLicenseRepo) ConfirmBlockchainStorage(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.QueryRowxContext(ctx, updateBlockchainConfirmationQuery,
		dl.BlockchainTxHash, dl.OnBlockchain, dl.ModifierId, dl.UpdatedAt, dl.Id,
	).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.ConfirmBlockchainStorage.StructScan")
	}
	return d, nil
}

func (r *DriverLicenseRepo) UpdateWalletAddress(ctx context.Context, dl *models.DrivingLicense) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.QueryRowxContext(ctx, updateWalletAddressQuery,
		dl.WalletAddress, dl.ModifierId, dl.UpdatedAt, dl.Id,
	).StructScan(d); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.UpdateWalletAddress.StructScan")
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

func (r *DriverLicenseRepo) GetDriverLicenseByWalletAddress(ctx context.Context, address string) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.GetContext(ctx, d, getDriverLicenseByWalletAddressQuery, address); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.getDriverLicenseByWalletAddressQuery.GetContext")
	}
	return d, nil
}

func (r *DriverLicenseRepo) GetDriverLicenseByLicenseNO(ctx context.Context, address string) (*models.DrivingLicense, error) {
	d := &models.DrivingLicense{}
	if err := r.db.GetContext(ctx, d, getDriverLicenseByLicenseNOQuery, address); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetDriverLicenseByLicenseNO.GetContext")
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

func (r *DriverLicenseRepo) GetStatusDistribution(ctx context.Context) (*models.StatusDistributionResponse, error) {
	var items []struct {
		Status string `db:"status"`
		Count  int    `db:"count"`
	}

	if err := r.db.SelectContext(ctx, &items, getStatusDistributionQuery); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetStatusDistribution.SelectContext")
	}

	total := 0
	distribution := make([]models.StatusDistributionItem, len(items))
	for i, item := range items {
		distribution[i] = models.StatusDistributionItem{
			Status: item.Status,
			Count:  item.Count,
		}
		total += item.Count
	}

	return &models.StatusDistributionResponse{
		Distribution: distribution,
		Total:        total,
	}, nil
}

func (r *DriverLicenseRepo) GetLicenseTypeDistribution(ctx context.Context) (*models.LicenseTypeDistributionResponse, error) {
	var items []struct {
		LicenseType string `db:"license_type"`
		Count       int    `db:"count"`
	}

	if err := r.db.SelectContext(ctx, &items, getLicenseTypeDistributionQuery); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetLicenseTypeDistribution.SelectContext")
	}

	total := 0
	distribution := make([]models.LicenseTypeDistributionItem, len(items))
	for i, item := range items {
		distribution[i] = models.LicenseTypeDistributionItem{
			LicenseType: item.LicenseType,
			Count:       item.Count,
		}
		total += item.Count
	}

	return &models.LicenseTypeDistributionResponse{
		Distribution: distribution,
		Total:        total,
	}, nil
}

func (r *DriverLicenseRepo) GetLicenseTypeStatusDistribution(ctx context.Context) (*models.LicenseTypeDetailDistributionResponse, error) {
	type row struct {
		LicenseType string `db:"license_type"`
		Status      string `db:"status"`
		Count       int    `db:"count"`
	}

	var rows []row
	if err := r.db.SelectContext(ctx, &rows, getLicenseTypeStatusDistributionQuery); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetLicenseTypeStatusDistribution.SelectContext")
	}

	// group data by license_type
	distMap := make(map[string]*models.LicenseTypeDetailDistribution)
	grandTotal := 0

	for _, r := range rows {
		grandTotal += r.Count

		if _, ok := distMap[r.LicenseType]; !ok {
			distMap[r.LicenseType] = &models.LicenseTypeDetailDistribution{
				LicenseType: r.LicenseType,
				Total:       0,
				ByStatus:    make([]models.StatusDistributionItem, 0),
			}
		}

		item := distMap[r.LicenseType]
		item.Total += r.Count
		item.ByStatus = append(item.ByStatus, models.StatusDistributionItem{
			Status: r.Status,
			Count:  r.Count,
		})
	}

	// change map to slice
	distribution := make([]models.LicenseTypeDetailDistribution, 0, len(distMap))
	for _, v := range distMap {
		distribution = append(distribution, *v)
	}

	// sort by total quantity in desc order
	sort.Slice(distribution, func(i, j int) bool {
		return distribution[i].Total > distribution[j].Total
	})

	return &models.LicenseTypeDetailDistributionResponse{
		Distribution: distribution,
		GrandTotal:   grandTotal,
	}, nil
}

func (r *DriverLicenseRepo) GetCityStatusDistribution(ctx context.Context) (*models.CityDetailDistributionResponse, error) {
	type row struct {
		OwnerCity string `db:"owner_city"`
		Status    string `db:"status"`
		Count     int    `db:"count"`
	}

	var rows []row
	if err := r.db.SelectContext(ctx, &rows, getCityStatusDistributionQuery); err != nil {
		return nil, errors.Wrap(err, "DriverLicenseRepo.GetCityStatusDistribution.SelectContext")
	}

	// group owner_city
	distMap := make(map[string]*models.CityDetailDistribution)
	grandTotal := 0

	for _, r := range rows {
		grandTotal += r.Count

		if _, ok := distMap[r.OwnerCity]; !ok {
			distMap[r.OwnerCity] = &models.CityDetailDistribution{
				OwnerCity: r.OwnerCity,
				Total:     0,
				ByStatus:  make([]models.CityStatusItem, 0),
			}
		}

		item := distMap[r.OwnerCity]
		item.Total += r.Count
		item.ByStatus = append(item.ByStatus, models.CityStatusItem{
			Status: r.Status,
			Count:  r.Count,
		})
	}

	distribution := make([]models.CityDetailDistribution, 0, len(distMap))
	for _, v := range distMap {
		distribution = append(distribution, *v)
	}

	// sort by total quantity in desc order
	sort.Slice(distribution, func(i, j int) bool {
		return distribution[i].Total > distribution[j].Total
	})

	return &models.CityDetailDistributionResponse{
		Distribution: distribution,
		GrandTotal:   grandTotal,
	}, nil
}
