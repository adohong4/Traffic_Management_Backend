package main

import (
	"log"
	"os"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/server"
	"github.com/adohong4/driving-license/pkg/db/postgres"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
)

// @title Driving License REST API
// @version 1.0
// @description REST API for Driving License Management
// @contact.url https://github.com/adohong4
// @BasePath /api/v1
func main() {
	log.Println("Starting driving license API server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	// Read & Analyst Config
	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	// Initialize Logger
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	// Connect PostgreSQL
	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %v", err)
	}
	defer psqlDB.Close()
	appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())

	// Run Server
	s := server.NewServer(cfg, psqlDB, appLogger)
	if err := s.Run(); err != nil {
		log.Fatalf("Server run: %v", err)
	}
}
