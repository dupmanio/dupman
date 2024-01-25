package errors

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type HTTPError struct {
	Code    int `json:"code"`
	Message any `json:"error,omitempty"`
}

func NewHTTPError(response *resty.Response) *HTTPError {
	if err, ok := response.Error().(*HTTPError); ok {
		return err
	}

	errorMessage := response.Error()
	if errorMessage == nil {
		errorMessage = response.String()
	}

	return &HTTPError{
		Code:    response.StatusCode(),
		Message: errorMessage,
	}
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("%s [%d]", err.Message, err.Code)
}
