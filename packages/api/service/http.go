package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/dto"
	"github.com/dupmanio/dupman/packages/dbutils/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HTTPService struct{}

func NewHTTPService() *HTTPService {
	return &HTTPService{}
}

func (svc *HTTPService) HTTPError(ctx *gin.Context, code int, err any) {
	ctx.AbortWithStatusJSON(code, dto.HTTPResponse{
		Code:  code,
		Error: err,
	})
}

func (svc *HTTPService) HTTPResponse(ctx *gin.Context, code int, data any) {
	ctx.JSON(code, dto.HTTPResponse{
		Code: code,
		Data: data,
	})
}

func (svc *HTTPService) HTTPPaginatedResponse(ctx *gin.Context, code int, data any, pagination *pagination.Pagination) {
	ctx.JSON(code, dto.HTTPResponse{
		Code:       code,
		Data:       data,
		Pagination: pagination,
	})
}

func (svc *HTTPService) HTTPValidationError(ctx *gin.Context, err error) {
	svc.HTTPError(ctx, http.StatusBadRequest, svc.normalizeHTTPValidationError(err))
}

func (svc *HTTPService) normalizeHTTPValidationError(err error) []string {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		return svc.formatValidationErrors(validationErr)
	}

	return []string{err.Error()}
}

func (svc *HTTPService) formatValidationErrors(validationErrors validator.ValidationErrors) []string {
	messages := make([]string, 0, len(validationErrors))

	for _, fieldError := range validationErrors {
		var errorMessage string

		switch fieldError.Tag() {
		case "required":
			errorMessage = fmt.Sprintf("Key '%s' is required", fieldError.Field())
		case "url":
			errorMessage = fmt.Sprintf("Value of field '%s' is not a valid URL address", fieldError.Field())
		default:
			errorMessage = fieldError.Error()
		}

		messages = append(messages, errorMessage)
	}

	return messages
}