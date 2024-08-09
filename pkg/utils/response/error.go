package response

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewError(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}
