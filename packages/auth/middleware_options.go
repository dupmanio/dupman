package auth

type Option func(*Middleware)

func WithHTTPErrorHandler(handler HTTPErrorHandlerFunc) Option {
	return func(mid *Middleware) {
		mid.httpErrorHandler = handler
	}
}

func WithCallUserService(callUserService bool) Option {
	return func(mid *Middleware) {
		mid.callUserService = callUserService
	}
}

func WithResetFilters() Option {
	return func(mid *Middleware) {
		mid.filters = make([]Filter, 0)
	}
}

func WithFilters(filters ...Filter) Option {
	return func(mid *Middleware) {
		mid.filters = append(mid.filters, filters...)
	}
}
