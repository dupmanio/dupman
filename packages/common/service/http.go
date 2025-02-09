package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type HTTPService struct {
	authSvc *AuthService
	ot      *otel.OTel
}

func NewHTTPService(authSvc *AuthService, ot *otel.OTel) *HTTPService {
	return &HTTPService{
		authSvc: authSvc,
		ot:      ot,
	}
}

func (svc *HTTPService) HTTPError(ctx *gin.Context, code int, err any) {
	ctx.AbortWithStatusJSON(code, dto.HTTPResponse[any]{
		Code:  code,
		Error: err,
	})
}

func (svc *HTTPService) HTTPErrorWithOTelLog(
	ctx *gin.Context,
	description string,
	code int,
	err error,
	attributes ...attribute.KeyValue,
) {
	svc.ot.LogErrorEvent(ctx, description, err, attributes...)
	svc.HTTPError(ctx, code, err.Error())
}

func (svc *HTTPService) HTTPValidationErrorWithOTelLog(
	ctx *gin.Context,
	err error,
	attributes ...attribute.KeyValue,
) {
	svc.ot.LogErrorEvent(ctx, "Invalid Payload", err, attributes...)
	svc.HTTPValidationError(ctx, err)
}

func (svc *HTTPService) HTTPResponse(ctx *gin.Context, code int, data any) {
	ctx.JSON(code, dto.HTTPResponse[any]{
		Code: code,
		Data: data,
	})
}

func (svc *HTTPService) SSEEvent(ctx *gin.Context, name string, message any) {
	ctx.SSEvent(name, message)
	ctx.Writer.Flush()
}

func (svc *HTTPService) HTTPPaginatedResponse(ctx *gin.Context, code int, data any, pagination *pagination.Pagination) {
	ctx.JSON(code, dto.HTTPResponse[any]{
		Code:       code,
		Data:       data,
		Pagination: pagination,
	})
}

func (svc *HTTPService) HTTPValidationError(ctx *gin.Context, err error) {
	svc.HTTPError(ctx, http.StatusBadRequest, svc.normalizeHTTPValidationError(err))
}

func (svc *HTTPService) EnrichSpanWithControllerAttributes(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	_, functionAttributes := otel.GetFunctionCallAttributes(1)

	span.SetAttributes(functionAttributes...)
	span.SetAttributes(semconv.EnduserID(svc.authSvc.CurrentUserID(ctx).String()))
}

func (svc *HTTPService) normalizeHTTPValidationError(err error) []string {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		return svc.formatValidationErrors(validationErr)
	}

	return []string{"invalid payload"}
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
		case "email":
			errorMessage = fmt.Sprintf("Value of field '%s' is not a valid Email address", fieldError.Field())
		default:
			errorMessage = fieldError.Error()
		}

		messages = append(messages, errorMessage)
	}

	return messages
}
