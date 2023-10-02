package constant

type ContextKey string

const (
	// CurrentUserKey represents key for storing Current User object.
	CurrentUserKey = "current_user"

	// EncryptionKeyKey represents key for website encryption key.
	EncryptionKeyKey ContextKey = "encryption_key"

	// PublicKeyHeaderKey represents key for the Public Key header.
	PublicKeyHeaderKey = "X-Public-Key"
)
