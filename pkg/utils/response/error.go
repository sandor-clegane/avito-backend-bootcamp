package response

type errorResponse struct {
	Error string `json:"error"`
}

func NewError(err error) errorResponse {
	return errorResponse{
		Error: err.Error(),
	}
}
