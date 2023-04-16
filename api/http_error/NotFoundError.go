package http_error

import "net/http"

func NotFoundError(entityName string, id string) HttpError {
	return HttpError{
		EntityName: entityName,
		Status:     http.StatusNotFound,
		ID:         id,
		Message:    "Not found " + entityName + " with identifier " + id,
	}
}
