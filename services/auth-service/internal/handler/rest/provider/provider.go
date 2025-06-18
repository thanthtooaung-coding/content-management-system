package provider

import (
	"go.uber.org/fx"
)

var Module = fx.Module("handler_module", fx.Provide())
