package http_error

type HttpError struct {
	ID         string `json:"id"`
	Status     int    `json:"status"`
	EntityName string `json:"entityName"`
	Message    string `json:"message"`
}

func (h HttpError) Error() string {
	return "Http Error"
}
