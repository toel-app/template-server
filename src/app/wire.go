//go:build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/toel-app/template-server/src/internal/ping"
	"github.com/toel-app/template-server/src/pkg/config"
	"github.com/toel-app/template-server/src/pkg/db"
	"github.com/toel-app/template-server/src/pkg/logger"
)

func Wire() Application {
	wire.Build(
		NewApp,
		NewRouter,
		logger.NewLogger,
		config.NewConfig,
		db.NewCloser,
		db.NewMongoDB,
		ping.Set,
	)

	return app{}
}
