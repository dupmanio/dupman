package dto

import "github.com/dupmanio/dupman/packages/dbutils/pagination"

type HTTPResponse[T any] struct {
	Code       int                    `json:"code"`
	Data       T                      `json:"data,omitempty"`
	Error      any                    `json:"error,omitempty"`
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
}
