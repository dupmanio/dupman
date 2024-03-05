package service

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type Option func(request *resty.Request)

func WithContext(ctx context.Context) Option {
	return func(request *resty.Request) {
		request.SetContext(ctx)
	}
}

func ApplyOptions(request *resty.Request, options []Option) {
	for _, opt := range options {
		opt(request)
	}
}
