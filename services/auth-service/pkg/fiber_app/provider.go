package fiber_app

import (
	"context"
	"fmt"
	"os"

	graph2 "github.com/content-management-system/auth-service/internal/handler/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/content-management-system/auth-service/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewFiberApp)

type FiberApp struct {
	App    *fiber.App
	logger *logrus.Logger
	db     *db.DB
}

func NewFiberApp(lifeCycle fx.Lifecycle, log *logrus.Logger, db *db.DB) *FiberApp {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	fiberApp := &FiberApp{
		App:    app,
		logger: log,
		db:     db,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fiberApp.setupRoutes()

	lifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fiberApp.logger.Info(fmt.Sprintf("Starting Fiber server on :%s", port))
			go func() {
				if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
					fiberApp.logger.WithError(err).Fatal("Failed to start Fiber server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fiberApp.logger.Info("Shutting down Fiber server")
			return app.Shutdown()
		},
	})

	return fiberApp
}

func (app *FiberApp) setupRoutes() {
	app.App.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Health check endpoint
	app.App.Get("/health", func(c *fiber.Ctx) error {
		var result int
		if err := app.db.Conn.Raw("SELECT 1").Scan(&result).Error; err != nil {
			app.logger.WithError(err).Error("Database health check failed")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"error":  "Database connection failed",
			})
		}
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Application is healthy",
		})
	})

}

func (app *FiberApp) setupGraphQL(resolver *graph2.Resolver) {
	srv := handler.NewDefaultServer(graph2.NewExecutableSchema(graph2.Config{Resolvers: resolver}))

	app.App.All("/graph", adaptor.HTTPHandler(srv))

	app.App.Get("/playground", adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/graph")))
}
