package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toel-app/common-utils/logger"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error(fmt.Sprintf("server panic %s", err), nil)
			c.JSON(http.StatusInternalServerError, NewError(ERR_INTERNAL_SERVER_ERROR))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
