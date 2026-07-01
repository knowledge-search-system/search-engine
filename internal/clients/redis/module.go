package redis

import "go.uber.org/fx"

var Module = fx.Module("redis_client",
	fx.Provide(NewClient),
)
