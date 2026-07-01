package server

import "go.uber.org/fx"

var Module = fx.Module("server",
	fx.Provide(
		newGRPCServer,
		newHTTPHandler,
	),
	fx.Invoke(
		registerGRPCLifecycle,
		registerHTTPLifecycle,
	),
)
