package http_error

import "net/http"

func Unauthorized(entityName string, id string) HttpError {
	return HttpError{
		ID:         id,
		EntityName: entityName,
		Status:     http.StatusUnauthorized,
		Message:    "Operation Not Allowed",
	}
}
