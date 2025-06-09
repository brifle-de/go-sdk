package api

type ResponseStatus struct {
	ErrorCode  int
	HttpStatus int
}

type ResponseError struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
