package ping

import (
	"github.com/toel-app/registration/src/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface{}

type controller struct {
}

func (r controller) ping(c *gin.Context) {
	c.JSON(http.StatusOK, response.NewOK("pong"))
}

func NewController(r *gin.Engine) Controller {
	res := controller{}
	r.GET("/ping", res.ping)

	return &res
}
