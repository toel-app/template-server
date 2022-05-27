package utils

import (
	"net/http"
)

// Register error here
const (
	ERR_INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	ERR_INVALID_PAYLOAD       = "INVALID_PAYLOAD"
	ERR_INVALID_OTP_METHOD    = "INVALID_OTP_METHOD"
	ERR_INVALID_OTP_CODE      = "INVALID_OTP_CODE"
)

// Don't forget to map error code here
var errorList = map[string]ErrResponse{
	ERR_INTERNAL_SERVER_ERROR: {
		Message:    "Internal server error",
		Code:       ERR_INTERNAL_SERVER_ERROR,
		HttpStatus: http.StatusInternalServerError,
	},
	ERR_INVALID_PAYLOAD: {
		Message:    "Invalid payload",
		Code:       ERR_INVALID_PAYLOAD,
		HttpStatus: http.StatusBadRequest,
	},
	ERR_INVALID_OTP_METHOD: {
		Message:    "Invalid OTP method",
		Code:       ERR_INVALID_OTP_METHOD,
		HttpStatus: http.StatusBadRequest,
	},
	ERR_INVALID_OTP_CODE: {
		Message:    "Invalid OTP code",
		Code:       ERR_INVALID_OTP_CODE,
		HttpStatus: http.StatusBadRequest,
	},
}

type RestResponse struct {
	Message string      `json:"message,omitempty"`
	Status  int         `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	Message    string `json:"message,omitempty"`
	Code       string `json:"code,omitempty"`
	HttpStatus int    `json:"-"`
}

func NewError(code string) *ErrResponse {
	for k, v := range errorList {
		if k == code {
			httpStatus := v.HttpStatus
			if httpStatus == 0 {
				httpStatus = http.StatusInternalServerError
			}
			return &ErrResponse{
				Message:    v.Message,
				Code:       v.Code,
				HttpStatus: v.HttpStatus,
			}
		}
	}

	return &ErrResponse{
		Message:    "Internal server error",
		Code:       ERR_INTERNAL_SERVER_ERROR,
		HttpStatus: http.StatusInternalServerError,
	}
}

func NewOK(data interface{}) *RestResponse {
	return &RestResponse{
		Data:   data,
		Status: http.StatusOK,
	}
}
