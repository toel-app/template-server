package util

import (
	"fmt"
	"github.com/toel-app/registration/src/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			fmt.Printf("server panic %s", err)
			c.JSON(http.StatusInternalServerError, response.NewError(response.ErrInternalServerError))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
