package response

type errorResponse struct {
	Error error `json:"error"`
}

func NewError(err error) errorResponse {
	return errorResponse{
		Error: err,
	}
}
