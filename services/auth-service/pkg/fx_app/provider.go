package fx_app

import (
	"github.com/content-management-system/auth-service/pkg/db"
	"github.com/content-management-system/auth-service/pkg/fiber_app"
	"github.com/sirupsen/logrus"
)

type App struct {
	DB     *db.DB
	Logger *logrus.Logger
	App    *fiber_app.FiberApp
}
