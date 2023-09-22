package structures

type CommonResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	CommonResponse
}

type SuccessResponse struct {
	CommonResponse
}
