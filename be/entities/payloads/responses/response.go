package responses

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewResponse(message string, data any, err error, status int) *Response {
	response := &Response{
		Status:  status,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	if data == nil {
		return response
	}
	response.Data = data
	return response
}
