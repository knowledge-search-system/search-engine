package searchhistory

import "go.uber.org/fx"

var Module = fx.Module("searchhistory_repository",
	fx.Provide(NewRepository),
)
