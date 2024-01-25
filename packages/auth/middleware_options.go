package auth

type Option func(*Middleware)

func WithHTTPErrorHandler(handler HTTPErrorHandlerFunc) Option {
	return func(mid *Middleware) {
		mid.httpErrorHandler = handler
	}
}

func WithFetchUserData(fetchUserData bool) Option {
	return func(mid *Middleware) {
		mid.fetchUserData = fetchUserData
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
