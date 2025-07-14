package usecase

import (
	"context"
	"testing"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth/mock"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUC_CreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	user := &models.User{
		IdentityNo: "123456789",
		Password:   "password123",
	}

	ctx := context.Background()

	// Mock FindByIdentity to return no existing user
	mockAuthRepo.EXPECT().FindByIdentity(ctx, gomock.Eq(user)).Return(nil, nil)
	// Mock CreateUser to return the created user
	mockAuthRepo.EXPECT().CreateUser(ctx, gomock.Eq(user)).Return(user, nil)

	createdUser, err := authUC.CreateUser(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.NotEmpty(t, createdUser.Token)
	require.Equal(t, user.IdentityNo, createdUser.User.IdentityNo)
}

func TestAuthUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	user := &models.User{
		Id:         uuid.New(),
		IdentityNo: "123456789",
		Password:   "password123",
	}

	ctx := context.Background()

	// Mock Update to return the updated user
	mockAuthRepo.EXPECT().Update(ctx, gomock.Eq(user)).Return(user, nil)

	updatedUser, err := authUC.Update(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Equal(t, user.Id, updatedUser.Id)
	require.Equal(t, user.IdentityNo, updatedUser.IdentityNo)
}

func TestAuthUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	userID := uuid.New()
	modifierID := uuid.New()
	version := 1

	ctx := context.Background()

	// Mock Delete to return no error
	mockAuthRepo.EXPECT().Delete(ctx, gomock.Eq(userID), gomock.Eq(modifierID), gomock.Eq(version)).Return(nil)

	err := authUC.Delete(ctx, userID, modifierID, version)
	require.NoError(t, err)
}

func TestAuthUC_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	userID := uuid.New()
	user := &models.User{
		Id:         userID,
		IdentityNo: "123456789",
		Password:   "password123",
	}

	ctx := context.Background()

	// Mock GetUserById to return the user
	mockAuthRepo.EXPECT().GetUserById(ctx, gomock.Eq(userID)).Return(user, nil)

	fetchedUser, err := authUC.GetByID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, fetchedUser)
	require.Equal(t, userID, fetchedUser.Id)
	require.Equal(t, user.IdentityNo, fetchedUser.IdentityNo)
}

func TestAuthUC_FindByIdentity(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	identity := "123456789"
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	ctx := context.Background()

	usersList := &models.UsersList{}

	// Mock FindByIdentityNO to return users list
	mockAuthRepo.EXPECT().FindByIdentityNO(ctx, gomock.Eq(identity), gomock.Eq(query)).Return(usersList, nil)

	userList, err := authUC.FindByIdentity(ctx, identity, query)
	require.NoError(t, err)
	require.NotNil(t, userList)
}

func TestAuthUC_GetUsers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}

	ctx := context.Background()

	usersList := &models.UsersList{}

	// Mock GetUsers to return users list
	mockAuthRepo.EXPECT().GetUsers(ctx, gomock.Eq(query)).Return(usersList, nil)

	users, err := authUC.GetUsers(ctx, query)
	require.NoError(t, err)
	require.NotNil(t, users)
}

func TestAuthUC_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, apiLogger)

	ctx := context.Background()

	user := &models.User{
		IdentityNo: "123456789",
		Password:   "password123",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Id:           uuid.New(),
		IdentityNo:   user.IdentityNo,
		HashPassword: string(hashPassword),
		Active:       true,
	}

	// Mock FindByIdentity to return the mock user
	mockAuthRepo.EXPECT().FindByIdentity(ctx, gomock.Any()).Return(mockUser, nil).Do(func(ctx context.Context, u *models.User) {
		t.Logf("FindByIdentity input: IdentityNO=%s, Password=%s", u.IdentityNo, u.Password)
		require.Equal(t, user.IdentityNo, u.IdentityNo)
		require.Equal(t, "password123", u.Password)
	})

	userWithToken, err := authUC.Login(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, userWithToken)
	require.NotEmpty(t, userWithToken.Token)
	require.Equal(t, mockUser.IdentityNo, userWithToken.User.IdentityNo)
}
