package http_error

import "net/http"

func ConflictError(entityName string, id string) HttpError {
	return HttpError{
		ID:         id,
		Status:     http.StatusConflict,
		EntityName: entityName,
		Message:    "There is an existing " + entityName + " with identifier " + id,
	}
}
