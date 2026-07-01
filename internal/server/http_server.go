package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/knowledge-search-system/search-engine/config"
	pkggrpc "github.com/knowledge-search-system/search-engine/pkg/grpc"
	searchenginev1 "github.com/knowledge-search-system/search-engine/proto/searchengine/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newHTTPHandler(cfg *config.Config) (http.Handler, error) {
	gwMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(header string) (string, bool) {
			if strings.EqualFold(header, pkggrpc.LangMetadataKey) {
				return pkggrpc.LangMetadataKey, true
			}
			return runtime.DefaultHeaderMatcher(header)
		}),
	)

	endpoint := fmt.Sprintf("localhost:%d", cfg.GRPC.Port)
	if err := searchenginev1.RegisterSearchServiceHandlerFromEndpoint(
		context.Background(),
		gwMux,
		endpoint,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return nil, fmt.Errorf("register grpc-gateway handler: %w", err)
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/docs", newSwaggerHandler())
	rootMux.Handle("/docs/", newSwaggerHandler())
	rootMux.Handle("/", gwMux)

	return rootMux, nil
}

func registerHTTPLifecycle(lc fx.Lifecycle, httpHandler http.Handler, cfg *config.Config, logger *zap.Logger) {
	srv := &http.Server{Handler: httpHandler}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.HTTP.Port)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("listen http on %s: %w", addr, err)
			}

			go func() {
				logger.Info("http server started", zap.String("addr", addr))
				if err := srv.Serve(lis); err != nil && err != http.ErrServerClosed {
					logger.Error("http server stopped", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
