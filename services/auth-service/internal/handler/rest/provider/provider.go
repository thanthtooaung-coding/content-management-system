package provider

import (
	"github.com/content-management-system/auth-service/internal/handler/rest"
	"go.uber.org/fx"
)

var Module = fx.Provide(rest.NewHandler)
