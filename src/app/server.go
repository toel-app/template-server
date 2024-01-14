package app

import (
	"github.com/toel-app/template-server/src/internal/ping"
	"github.com/toel-app/template-server/src/pkg/db"
	log "github.com/toel-app/template-server/src/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

type Application interface {
	Run()
	Shutdown()
}

type app struct {
	logger      log.Logger
	dbCloser    db.Closer
	pingHandler ping.Controller
}

func NewApp(
	logger log.Logger,
	dbCloser db.Closer,
	pingHandler ping.Controller,
) Application {
	return &app{
		logger,
		dbCloser,
		pingHandler,
	}
}

// Shutdown do closing infrastructural things e.g closing db
func (a app) Shutdown() {
	if err := a.dbCloser.Close(); err != nil {
		a.logger.Error("cannot close db", err)
		return
	}
}

func (a app) Run() {
	a.logger.Info("server started at port:8080")
	if err := router.Run(); err != nil {
		a.logger.Error("server not started", err)
	}
}

func StartApplication() {
	application := Wire()
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go application.Run()
	<-cancelChan
	application.Shutdown()
}
