package response

type Responses struct {
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Status    any         `json:"status"`
	Meta_data MetaData    `json:"meta_data"`
}

type ResponsesOne struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  any         `json:"status"`
}

type FailedResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Status  any         `json:"status"`
	Error   any         `json:"error"`
}

type MetaData struct {
	Limit    int    `json:"limit"`
	Pages    int    `json:"pages"`
	Total    int    `json:"total"`
	Sort_by  string `json:"sort_by"`
	Sort_key string `json:"sort_key"`
}

func NewResponse(data interface{}, message string, status any) Responses {
	response := Responses{
		Data:    data,
		Message: message,
		Status:  status,
	}

	return response
}

func OneResponse(data interface{}, message string, status any) ResponsesOne {
	response := ResponsesOne{
		Data:    data,
		Message: message,
		Status:  status,
	}

	return response
}

func NewFailedResponse(data interface{}, message string, status any, err any) FailedResponse {
	response := FailedResponse{
		Data:    data,
		Message: message,
		Status:  status,
		Error:   err,
	}

	return response
}
