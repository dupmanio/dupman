package constant

type ContextKey string

const (
	// UserIDKey represents key for storing authenticated user ID.
	UserIDKey = "user_id"

	// EncryptionKeyKey represents key for website encryption key.
	EncryptionKeyKey ContextKey = "encryption_key"
)
