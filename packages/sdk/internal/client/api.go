package client

import (
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/go-resty/resty/v2"
)

func NewAPIClient(sess *session.Session) *resty.Client {
	// @todo: set url by ENV.
	return getBaseClient(sess.Config, sess.Token.AccessToken).
		SetBaseURL("http://gateway.dupman.localhost/api")
}
