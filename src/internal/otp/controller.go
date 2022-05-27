package otp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toel-app/common-utils/logger"
	"github.com/toel-app/common-utils/string_utils"
	"github.com/toel-app/template-server/src/pkg/utils"
)

func RegisterRoute(r *gin.Engine, service IService) {
	res := resource{service}
	r.Group("/otp").
		POST("send", res.Send).
		POST("verify", res.VerifyOtp)
}

type resource struct {
	service IService
}

func (r resource) Send(c *gin.Context) {
	var payload CreateOtp
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error("error binding payload", err)
		c.JSON(http.StatusBadRequest, utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR))
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(err.HttpStatus, utils.NewError(err.Code))
		return
	}

	_, isEmail := string_utils.ValidMailAddress(payload.Addressee)
	if !isEmail {
		payload.Addressee = string_utils.FormatPhoneNumber(payload.Addressee)
	}
	if err := r.service.Send(payload); err != nil {
		c.JSON(err.HttpStatus, utils.NewError(err.Code))
		return
	}

	c.JSON(http.StatusCreated, utils.NewOK("OTP sent"))
}

func (r resource) VerifyOtp(c *gin.Context) {
	var payload VerifyOtp
	if err := c.ShouldBindJSON(&payload); err != nil {
		logger.Error("error binding payload", err)
		c.JSON(http.StatusBadRequest, utils.NewError(utils.ERR_INTERNAL_SERVER_ERROR))
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(err.HttpStatus, utils.NewError(err.Code))
		return
	}

	err := r.service.VerifyOtp(payload)

	if err != nil {
		c.JSON(err.HttpStatus, utils.NewError(err.Code))
		return
	}

	c.JSON(http.StatusCreated, utils.NewOK("OTP verified!"))
}
