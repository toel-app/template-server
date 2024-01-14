package response

import "errors"

var (
	ErrInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
	ErrResourceNotFound    = errors.New("RESOURCE_NOT_FOUND")
	ErrUnauthorized        = errors.New("UNAUTHORIZED")
)
