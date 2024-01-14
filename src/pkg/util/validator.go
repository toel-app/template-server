package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/toel-app/registration/src/pkg/logger"
	"net/http"
)

type ErrResponse struct {
	Message    string `json:"message,omitempty"`
	Code       string `json:"code,omitempty"`
	HttpStatus int    `json:"-"`
}

var errResponse = ErrResponse{
	Message:    "Invalid payload",
	Code:       "INVALID_PAYLOAD",
	HttpStatus: http.StatusBadRequest,
}

var log = logger.NewLogger()

func validateBodyWrapper[Payload any](p Payload, c *gin.Context) *ErrResponse {
	errBindingPayloadLog := fmt.Sprintf("error binding payload for path %s", c.Request.URL.Path)
	errValidateBody := fmt.Sprintf("error validate payload for path %s", c.Request.URL.Path)

	if err := c.ShouldBindBodyWith(&p, binding.JSON); err != nil {
		log.Error(errBindingPayloadLog, err)
		return &errResponse
	}

	err := validator.New().Struct(p)
	if err != nil {
		log.Error(errValidateBody, err)
		return &errResponse
	}

	return nil
}

func validateQueryWrapper[Query any](q Query, c *gin.Context) *ErrResponse {
	errBindingPayloadLog := fmt.Sprintf("error binding query for path %s", c.Request.URL.Path)
	errValidateQuery := fmt.Sprintf("error validate query for path %s", c.Request.URL.Path)

	if err := c.ShouldBindQuery(&q); err != nil {
		log.Error(errBindingPayloadLog, err)
		return &errResponse
	}

	err := validator.New().Struct(q)
	if err != nil {
		log.Error(errValidateQuery, err)
		return &errResponse
	}

	return nil
}

func ValidateBody[T any](c *gin.Context) {
	var payload T

	err := validateBodyWrapper(payload, c)

	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse)
		c.Abort()
		return
	}

	c.Next()
}

func ValidateQuery[T any](c *gin.Context) {
	var payload T

	err := validateQueryWrapper(payload, c)

	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse)
		c.Abort()
		return
	}

	c.Next()
}
