package response

import (
	"errors"
	"fmt"
	"net/http"
)

// Don't forget to map error code here
var errorList = map[error]Err{
	ErrInternalServerError: {
		Message:    "Internal server error",
		HttpStatus: http.StatusInternalServerError,
	},
	ErrResourceNotFound: {
		Message:    "Resource not found",
		HttpStatus: http.StatusNotFound,
	},
	ErrUnauthorized: {
		Message:    "Unauthorized",
		HttpStatus: http.StatusUnauthorized,
	},
}

type RestResponse struct {
	Message string      `json:"message,omitempty"`
	Status  int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Err struct {
	Message    string `json:"message,omitempty"`
	Code       string `json:"code,omitempty"`
	HttpStatus int    `json:"-"`
}

func (e Err) Error() string {
	return fmt.Sprintf("message: %s, code: %s, httpStatus: %s", e.Message, e.Code, e.HttpStatus)
}

func NewError(code error) *Err {
	for e, v := range errorList {
		if errors.Is(e, code) {
			httpStatus := v.HttpStatus
			if httpStatus == 0 {
				httpStatus = http.StatusInternalServerError
			}
			return &Err{
				Message:    v.Message,
				Code:       e.Error(),
				HttpStatus: v.HttpStatus,
			}
		}
	}

	return &Err{
		Message:    "Internal server error",
		Code:       ErrInternalServerError.Error(),
		HttpStatus: http.StatusInternalServerError,
	}
}

func NewOK(data interface{}) *RestResponse {
	return &RestResponse{
		Data: data,
	}
}
