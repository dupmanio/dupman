package encryptor

type Encryptor interface {
	SetPrivateKey(string) error
	PrivateKey() string
	SetPublicKey(string) error
	PublicKey() (string, error)
	GenerateKeyPair() error
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}
