package main

import (
	"github.com/knowledge-search-system/search-engine/config"
	"github.com/knowledge-search-system/search-engine/internal/apperrors"
	esclient "github.com/knowledge-search-system/search-engine/internal/clients/elasticsearch"
	redisclient "github.com/knowledge-search-system/search-engine/internal/clients/redis"
	"github.com/knowledge-search-system/search-engine/internal/handler"
	"github.com/knowledge-search-system/search-engine/internal/logger"
	"github.com/knowledge-search-system/search-engine/internal/repository"
	"github.com/knowledge-search-system/search-engine/internal/repository/searchhistory"
	"github.com/knowledge-search-system/search-engine/internal/server"
	"github.com/knowledge-search-system/search-engine/internal/service"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		config.Module,
		logger.Module,
		apperrors.Module,
		repository.Module,
		searchhistory.Module,
		esclient.Module,
		redisclient.Module,
		service.Module,
		handler.Module,
		server.Module,
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
