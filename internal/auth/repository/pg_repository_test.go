package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestAuthRepo_Register kiểm tra hàm Register
func TestAuthRepo_Register(t *testing.T) {
	t.Parallel()

	// Khởi tạo mock database
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Register Success", func(t *testing.T) {
		// Trường hợp thành công: Tạo người dùng mới
		role := "admin"
		creatorID := uuid.New()
		modifierID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		user := &models.User{
			Id:         uuid.New(),
			IdentityNo: "123456789",
			Password:   string(hashedPassword),
			Active:     true,
			Role:       &role,
			Version:    1,
			CreatorId:  &creatorID,
			ModifierId: modifierID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		rows := sqlmock.NewRows([]string{
			"id", "identity_no", "password", "active", "role", "version", "creator_id", "modifier_id", "created_at", "updated_at",
		}).AddRow(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role, user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		)

		mock.ExpectQuery(createUserQuery).WithArgs(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role,
			user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		).WillReturnRows(rows)

		createdUser, err := authRepo.Register(context.Background(), user)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, user.Id, createdUser.Id)
		require.Equal(t, user.IdentityNo, createdUser.IdentityNo)
	})

	t.Run("Register Error", func(t *testing.T) {
		user := &models.User{
			Id:         uuid.New(),
			IdentityNo: "123456789",
			Password:   "password123",
			Active:     true,
			Version:    1,
		}

		mock.ExpectQuery(createUserQuery).WithArgs(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role,
			user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		).WillReturnError(sql.ErrConnDone)

		createdUser, err := authRepo.Register(context.Background(), user)
		require.Error(t, err)
		require.Nil(t, createdUser)
		require.True(t, errors.Is(errors.Cause(err), sql.ErrConnDone)) // Sử dụng errors.Cause
	})
}

// TestAuthRepo_Update kiểm tra hàm Update
func TestAuthRepo_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Update Success", func(t *testing.T) {
		// Trường hợp thành công: Cập nhật người dùng
		role := "user"
		creatorID := uuid.New()
		modifierID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("newpassword"), bcrypt.DefaultCost)
		user := &models.User{
			Id:         uuid.New(),
			IdentityNo: "987654321",
			Password:   string(hashedPassword),
			Active:     true,
			Role:       &role,
			Version:    2,
			CreatorId:  &creatorID,
			ModifierId: modifierID,
			CreatedAt:  time.Now().Add(-time.Hour),
			UpdatedAt:  time.Now(),
		}

		rows := sqlmock.NewRows([]string{
			"id", "identity_no", "password", "active", "role", "version", "creator_id", "modifier_id", "created_at", "updated_at",
		}).AddRow(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role, user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		)

		mock.ExpectQuery(updateUserQuery).WithArgs(
			user.IdentityNo, user.Password, user.Active, user.Role,
			user.CreatorId, user.ModifierId, user.Id, user.Version-1,
		).WillReturnRows(rows)

		updatedUser, err := authRepo.Update(context.Background(), user)
		require.NoError(t, err)
		require.NotNil(t, updatedUser)
		require.Equal(t, user.Id, updatedUser.Id)
		require.Equal(t, user.IdentityNo, updatedUser.IdentityNo)
	})

	t.Run("Update No Rows", func(t *testing.T) {
		// Trường hợp lỗi: Không tìm thấy người dùng để cập nhật
		user := &models.User{
			Id:         uuid.New(),
			IdentityNo: "987654321",
			Password:   "newpassword",
			Active:     true,
			Version:    2,
		}

		mock.ExpectQuery(updateUserQuery).WithArgs(
			user.IdentityNo, user.Password, user.Active, user.Role,
			user.CreatorId, user.ModifierId, user.Id, user.Version,
		).WillReturnError(sql.ErrNoRows)

		updatedUser, err := authRepo.Update(context.Background(), user)
		require.Error(t, err)
		require.Nil(t, updatedUser)
		require.True(t, errors.Is(errors.Cause(err), sql.ErrNoRows))
	})
}

// TestAuthRepo_Delete kiểm tra hàm Delete
func TestAuthRepo_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("Delete Success", func(t *testing.T) {
		// Trường hợp thành công: Xóa người dùng (đặt active = false)
		uid := uuid.New()
		modifierID := uuid.New()
		version := 1

		mock.ExpectExec(deleteUserQuery).WithArgs(
			modifierID, sqlmock.AnyArg(), uid, version,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := authRepo.Delete(context.Background(), uid, modifierID, version)
		require.NoError(t, err)
	})

	t.Run("Delete No Rows", func(t *testing.T) {
		// Trường hợp lỗi: Không tìm thấy người dùng để xóa
		uid := uuid.New()
		modifierID := uuid.New()
		version := 1

		mock.ExpectExec(deleteUserQuery).WithArgs(
			modifierID, sqlmock.AnyArg(), uid, version,
		).WillReturnResult(sqlmock.NewResult(1, 0))

		err := authRepo.Delete(context.Background(), uid, modifierID, version)
		require.Error(t, err)
		require.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

// TestAuthRepo_GetUserById kiểm tra hàm GetUserById
func TestAuthRepo_GetUserById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("GetUserById Success", func(t *testing.T) {
		// Trường hợp thành công: Tìm thấy người dùng theo ID
		uid := uuid.New()
		role := "user"
		creatorID := uuid.New()
		modifierID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		user := &models.User{
			Id:         uid,
			IdentityNo: "123456789",
			Password:   string(hashedPassword),
			Active:     true,
			Role:       &role,
			Version:    1,
			CreatorId:  &creatorID,
			ModifierId: modifierID,
			CreatedAt:  time.Now().Add(-time.Hour),
			UpdatedAt:  time.Now(),
		}

		rows := sqlmock.NewRows([]string{
			"id", "identity_no", "password", "active", "role", "version", "creator_id", "modifier_id", "created_at", "updated_at",
		}).AddRow(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role, user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		)

		mock.ExpectQuery(getUserQuery).WithArgs(uid).WillReturnRows(rows)

		foundUser, err := authRepo.GetUserById(context.Background(), uid)
		require.NoError(t, err)
		require.NotNil(t, foundUser)
		require.Equal(t, user.Id, foundUser.Id)
		require.Equal(t, user.IdentityNo, foundUser.IdentityNo)
	})

	t.Run("GetUserById Not Found", func(t *testing.T) {
		// Trường hợp lỗi: Không tìm thấy người dùng
		uid := uuid.New()

		mock.ExpectQuery(getUserQuery).WithArgs(uid).WillReturnError(sql.ErrNoRows)

		foundUser, err := authRepo.GetUserById(context.Background(), uid)
		require.Error(t, err)
		require.Nil(t, foundUser)
		require.True(t, errors.Is(err, sql.ErrNoRows))
	})
}

// TestAuthRepo_FindByIdentity kiểm tra hàm FindByIdentity
func TestAuthRepo_FindByIdentity(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByIdentity Success", func(t *testing.T) {
		// Trường hợp thành công: Tìm thấy người dùng theo IdentityNo
		uid := uuid.New()
		role := "user"
		creatorID := uuid.New()
		modifierID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		user := &models.User{
			Id:         uid,
			IdentityNo: "123456789",
			Password:   string(hashedPassword),
			Active:     true,
			Role:       &role,
			Version:    1,
			CreatorId:  &creatorID,
			ModifierId: modifierID,
			CreatedAt:  time.Now().Add(-time.Hour),
			UpdatedAt:  time.Now(),
		}

		rows := sqlmock.NewRows([]string{
			"id", "identity_no", "password", "active", "role", "version", "creator_id", "modifier_id", "created_at", "updated_at",
		}).AddRow(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role, user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		)

		mock.ExpectQuery(findUserByIdentity).WithArgs(user.IdentityNo).WillReturnRows(rows)

		foundUser, err := authRepo.FindByIdentity(context.Background(), user)
		require.NoError(t, err)
		require.NotNil(t, foundUser)
		require.Equal(t, user.Id, foundUser.Id)
		require.Equal(t, user.IdentityNo, foundUser.IdentityNo)
	})

	t.Run("FindByIdentity Not Found", func(t *testing.T) {
		// Trường hợp không tìm thấy người dùng
		user := &models.User{
			IdentityNo: "123456789",
		}

		mock.ExpectQuery(findUserByIdentity).WithArgs(user.IdentityNo).WillReturnError(sql.ErrNoRows)

		foundUser, err := authRepo.FindByIdentity(context.Background(), user)
		require.NoError(t, err)
		require.Nil(t, foundUser)
	})
}

// TestAuthRepo_GetUsers kiểm tra hàm GetUsers
func TestAuthRepo_GetUsers(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("GetUsers Success", func(t *testing.T) {
		// Trường hợp thành công: Lấy danh sách người dùng
		uid := uuid.New()
		role := "user"
		creatorID := uuid.New()
		modifierID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		user := &models.User{
			Id:         uid,
			IdentityNo: "123456789",
			Password:   string(hashedPassword),
			Active:     true,
			Role:       &role,
			Version:    1,
			CreatorId:  &creatorID,
			ModifierId: modifierID,
			CreatedAt:  time.Now().Add(-time.Hour),
			UpdatedAt:  time.Now(),
		}

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		rows := sqlmock.NewRows([]string{
			"id", "identity_no", "password", "active", "role", "version", "creator_id", "modifier_id", "created_at", "updated_at",
		}).AddRow(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role, user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		)

		pq := &utils.PaginationQuery{
			Page:    1,
			Size:    10,
			OrderBy: "",
		}

		mock.ExpectQuery(getTotal).WillReturnRows(totalCountRows)
		mock.ExpectQuery(getUsers).WithArgs(pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()).WillReturnRows(rows)

		usersList, err := authRepo.GetUsers(context.Background(), pq)
		require.NoError(t, err)
		require.NotNil(t, usersList)
		require.Equal(t, 1, usersList.TotalCount)
		require.Len(t, usersList.Users, 1)
		require.Equal(t, user.Id, usersList.Users[0].Id)
	})

	t.Run("GetUsers Empty", func(t *testing.T) {
		// Trường hợp danh sách rỗng
		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		pq := &utils.PaginationQuery{
			Page:    1,
			Size:    10,
			OrderBy: "",
		}

		mock.ExpectQuery(getTotal).WillReturnRows(totalCountRows)

		usersList, err := authRepo.GetUsers(context.Background(), pq)
		require.NoError(t, err)
		require.NotNil(t, usersList)
		require.Equal(t, 0, usersList.TotalCount)
		require.Empty(t, usersList.Users)
	})
}

// TestAuthRepo_FindByIdentityNO kiểm tra hàm FindByIdentityNO
func TestAuthRepo_FindByIdentityNO(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	authRepo := NewAuthRepository(sqlxDB)

	t.Run("FindByIdentityNO Success", func(t *testing.T) {
		// Trường hợp thành công: Tìm thấy người dùng theo IdentityNo
		uid := uuid.New()
		role := "user"
		creatorID := uuid.New()
		modifierID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		identity := "123456789"

		user := &models.User{
			Id:         uid,
			IdentityNo: identity,
			Password:   string(hashedPassword),
			Active:     true,
			Role:       &role,
			Version:    1,
			CreatorId:  &creatorID,
			ModifierId: modifierID,
			CreatedAt:  time.Now().Add(-time.Hour),
			UpdatedAt:  time.Now(),
		}

		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		rows := sqlmock.NewRows([]string{
			"id", "identity_no", "password", "active", "role", "version", "creator_id", "modifier_id", "created_at", "updated_at",
		}).AddRow(
			user.Id, user.IdentityNo, user.Password, user.Active, user.Role, user.Version, user.CreatorId, user.ModifierId, user.CreatedAt, user.UpdatedAt,
		)

		pq := &utils.PaginationQuery{
			Page:    1,
			Size:    10,
			OrderBy: "",
		}

		mock.ExpectQuery(getTotalCount).WithArgs(identity).WillReturnRows(totalCountRows)
		mock.ExpectQuery(findUsers).WithArgs(identity, pq.GetOffset(), pq.GetLimit()).WillReturnRows(rows)

		usersList, err := authRepo.FindByIdentityNO(context.Background(), identity, pq)
		require.NoError(t, err)
		require.NotNil(t, usersList)
		require.Equal(t, 1, usersList.TotalCount)
		require.Len(t, usersList.Users, 1)
		require.Equal(t, user.Id, usersList.Users[0].Id)
	})

	t.Run("FindByIdentityNO Empty", func(t *testing.T) {
		// Trường hợp không tìm thấy người dùng
		identity := "123456789"
		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		pq := &utils.PaginationQuery{
			Page:    1,
			Size:    10,
			OrderBy: "",
		}

		mock.ExpectQuery(getTotalCount).WithArgs(identity).WillReturnRows(totalCountRows)

		usersList, err := authRepo.FindByIdentityNO(context.Background(), identity, pq)
		require.NoError(t, err)
		require.NotNil(t, usersList)
		require.Equal(t, 0, usersList.TotalCount)
		require.Empty(t, usersList.Users)
	})

	t.Run("FindByIdentityNO Empty String", func(t *testing.T) {
		// Trường hợp chuỗi identity rỗng
		identity := ""
		totalCountRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		pq := &utils.PaginationQuery{
			Page:    1,
			Size:    10,
			OrderBy: "",
		}

		mock.ExpectQuery(getTotalCount).WithArgs(identity).WillReturnRows(totalCountRows)

		usersList, err := authRepo.FindByIdentityNO(context.Background(), identity, pq)
		require.NoError(t, err)
		require.NotNil(t, usersList)
		require.Equal(t, 0, usersList.TotalCount)
		require.Empty(t, usersList.Users)
	})
}
