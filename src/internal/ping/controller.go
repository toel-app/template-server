package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toel-app/template-server/src/pkg/utils"
)

func RegisterRoute(r *gin.Engine) {
	res := resource{}
	r.GET("/ping", res.ping)
}

type resource struct{}

func (r resource) ping(c *gin.Context) {
	c.JSON(http.StatusOK, utils.NewOK("pong"))
}
