package middleware

import (
	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/pkg/logger"
)

type MiddlewareManager struct {
	authUC  auth.UseCase
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

func NewMiddlewareManager(authUC auth.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{authUC: authUC, cfg: cfg, origins: origins, logger: logger}
}
