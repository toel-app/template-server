package app

import (
	"github.com/gin-gonic/gin"
	"github.com/toel-app/template-server/src/pkg/util"
	"sync"
)

var (
	router *gin.Engine
	once   sync.Once
)

func NewRouter() *gin.Engine {
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(util.Recovery())
		router = r
	})

	return router
}
