package server

import (
	"context"
	"fmt"
	"net"

	"github.com/knowledge-search-system/search-engine/config"
	"github.com/knowledge-search-system/search-engine/internal/apperrors"
	"github.com/knowledge-search-system/search-engine/internal/handler"
	pkggrpc "github.com/knowledge-search-system/search-engine/pkg/grpc"
	searchenginev1 "github.com/knowledge-search-system/search-engine/proto/searchengine/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newGRPCServer(searchHandler *handler.SearchHandler, logger *zap.Logger, translator *apperrors.Translator) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			pkggrpc.RecoveryInterceptor(logger),
			pkggrpc.LoggingInterceptor(logger),
			pkggrpc.ErrorTranslationInterceptor(translator),
		),
	)

	searchenginev1.RegisterSearchServiceServer(srv, searchHandler)
	reflection.Register(srv)

	return srv
}

func registerGRPCLifecycle(lc fx.Lifecycle, srv *grpc.Server, cfg *config.Config, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.GRPC.Port)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("listen grpc on %s: %w", addr, err)
			}

			go func() {
				logger.Info("grpc server started", zap.String("addr", addr))
				if err := srv.Serve(lis); err != nil {
					logger.Error("grpc server stopped", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})
}
