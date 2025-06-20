package main

import (
	"context"
	"github.com/content-management-system/auth-service/internal/service"
	"log"

	"github.com/content-management-system/auth-service/internal/handler/rest/provider"
	"github.com/content-management-system/auth-service/pkg/db"
	"github.com/content-management-system/auth-service/pkg/fiber_app"
	"github.com/content-management-system/auth-service/pkg/fx_app"
	"github.com/content-management-system/auth-service/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func NewApp(
	db *db.DB,
	logger *logrus.Logger,

	fiberApp *fiber_app.FiberApp,
) *fx_app.App {
	return &fx_app.App{
		DB:     db,
		Logger: logger,
		App:    fiberApp,
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	app := fx.New(
		fx.Provide(logger.NewLogger),
		db.Module,
		provider.Module,
		service.Module,
		fiber_app.Module,
		fx.Provide(NewApp),
		fx.Invoke(func(app *fx_app.App) {
			app.Logger.Info("Application initialized")
			var result int
			if err := app.DB.Conn.Raw("SELECT 1").Scan(&result).Error; err != nil {
				app.Logger.WithError(err).Error("Failed to verify database connection")
			} else {
				app.Logger.Info("Database connection verified, result:", result)
			}
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	<-app.Done()

	if err := app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
