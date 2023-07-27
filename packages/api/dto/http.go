package dto

type HTTPResponse struct {
	Code       int `json:"code"`
	Data       any `json:"data,omitempty"`
	Error      any `json:"error,omitempty"`
	Pagination any `json:"pagination,omitempty"`
}
