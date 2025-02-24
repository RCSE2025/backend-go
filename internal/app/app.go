package app

import (
	"fmt"
	"github.com/RCSE2025/backend-go/internal/config"
	"github.com/RCSE2025/backend-go/internal/http/handlers"
	"github.com/RCSE2025/backend-go/internal/model"
	"github.com/RCSE2025/backend-go/pkg/httpserver"
	"github.com/RCSE2025/backend-go/pkg/logger"
	"github.com/RCSE2025/backend-go/pkg/logger/sl"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run_() {
	cfg := config.Get()
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	if err := model.RunMigrations(db); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Миграции применились")
}

func Run() {
	cfg := config.Get()

	log := logger.NewLogger(cfg.Production)
	//ctx, cancel := context.WithCancel(context.Background())

	log.Info("app config", slog.Any("config", cfg))

	log.Info(
		"starting user-api",
		slog.Bool("PRODUCTION", cfg.Production),
		slog.String("version", cfg.Version),
	)
	log.Debug("debug messages are enabled")

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})

	if err != nil {
		log.Error("error connecting to database", sl.Err(err))
		return
	}

	if err := model.RunMigrations(db); err != nil {
		log.Error("error applying migrations", sl.Err(err))
		return
	}

	fmt.Println("Migrations applied")

	//r := chi.NewRouter()
	if cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	//userRepo := repo.New(postgresDB)

	handlers.NewRouter(r, log)

	httpServer := httpserver.New(r, httpserver.Port(cfg.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: ", slog.Any("signal", s.String()))
	case err := <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify: %w", sl.Err(err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown: %w", sl.Err(err))
	}

	sqlDB, err := db.DB()
	if err != nil {

		log.Error("app - Run - db.DB", sl.Err(err))
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Error("app - Run - sqlDB.Close", sl.Err(err))
		return
	}
	log.Info("app - Run - db closed")

	log.Info("app - Run - exiting")
	//cancel()
}
