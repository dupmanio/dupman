package constant

type ContextKey string

const (
	// CurrentUserKey represents key for storing Current User object.
	CurrentUserKey = "current_user"

	// TokenScopesKey represents key for storing Current JWT Tokens Scopes.
	TokenScopesKey = "token_scopes"

	// EncryptionKeyKey represents key for website encryption key.
	EncryptionKeyKey ContextKey = "encryption_key"
)
