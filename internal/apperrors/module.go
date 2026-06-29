package apperrors

import "go.uber.org/fx"

var Module = fx.Module("apperrors",
	fx.Provide(NewTranslator),
)
