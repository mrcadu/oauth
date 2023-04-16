package http_error

import "net/http"

func BadRequest() HttpError {
	return HttpError{
		Status:  http.StatusBadRequest,
		Message: "Bad Request",
	}
}
