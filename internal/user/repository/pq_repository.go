package repository

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	user "github.com/adohong4/driving-license/internal/user"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) user.Repository {
	return &userRepo{db: db}
}

func (r *notificationRepo) CreateUser(ctx context.Context, db *models.User) (*models.User, error) {
}
