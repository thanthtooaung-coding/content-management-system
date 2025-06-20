package provider

import (
	authHandle "github.com/content-management-system/auth-service/internal/handler/rest/handler"
	"go.uber.org/fx"
)

var Module = fx.Module("handler_module", fx.Provide(
	authHandle.NewAuthHandler,
))
