package errors

import "errors"

var (
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrUserDoesNotExist            = errors.New("user does not exist")
	ErrWebsiteNotFound             = errors.New("website not found")
	ErrAuthorizationRequired       = errors.New("authentication required")
	ErrAccessIsForbidden           = errors.New("access is forbidden")
	ErrMissingScopes               = errors.New("missing scopes")
	ErrUnableToFetchUpdates        = errors.New("unable to fetch Website Updates")
	ErrNoDupmanEndpoint            = errors.New("/dupman endpoint does not exist")
	ErrDupmanEndpointAccessDenied  = errors.New("access denied on /dupman endpoint")
	ErrTokenIsExpired              = errors.New("token is expired")
	ErrUnsupportedNotificationType = errors.New("unsupported notification type")
	ErrInvalidUserID               = errors.New("user id is invalid or missing")
	ErrUnableToPublishMessage      = errors.New("unable to publish message")
)
