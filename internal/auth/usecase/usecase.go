package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-auth"
	cacheDuration = 3600
)

type authUC struct {
	cfg      *config.Config
	authRepo auth.Repository
	logger   logger.Logger
}

// Auth Usecase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, log logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, logger: log}
}

func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	existsUser, err := u.authRepo.FindByIdentity(ctx, user)
	if existsUser != nil || err != nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrIdentityAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// update existing user
func (u *authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Update.PrepareUpdate"))
	}

	updatedUser, err := u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()
	return updatedUser, nil
}

// delete user
func (u *authUC) Delete(ctx context.Context, Id uuid.UUID) error {
	if err := u.authRepo.Delete(ctx, Id); err != nil {
		return err
	}

	return nil
}

// Get user by id
func (u *authUC) GetByID(ctx context.Context, Id uuid.UUID) (*models.User, error) {
	user, err := u.authRepo.GetUserById(ctx, Id)
	if err != nil {
		return nil, err
	}

	user.SanitizePassword()

	return user, nil
}

// Find users by identityNO
func (u *authUC) FindByIdentity(ctx context.Context, identity string, query *utils.PaginationQuery) (*models.UsersList, error) {
	return u.authRepo.FindByIdentityNO(ctx, &identity, query)
}

// Get users with pagination
func (u *authUC) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	return u.authRepo.GetUsers(ctx, pq)
}

// Login user, return user model with jwt token
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.authRepo.FindByIdentity(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.Login.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Login.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Generate User Key
func (u *authUC) GenerateUserKey(Id string) string {
	return fmt.Sprintf("%s: %s", basePrefix, Id)
}
