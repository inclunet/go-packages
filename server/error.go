package server

type ErrorBody struct {
	Message string `json:"message"`
}

func NewError(statusCode int, err error) *Response {
	return &Response{
		StatusCode: statusCode,
		Body: &ErrorBody{
			Message: err.Error(),
		},
	}
}
