package encryptor

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type RSAEncryptor struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

const keySize = 2048

var errUnableToDecodeKey = errors.New("unable to decode key")

func NewRSAEncryptor() *RSAEncryptor {
	return &RSAEncryptor{}
}

// SetPrivateKey sets the Private Key value.
func (enc *RSAEncryptor) SetPrivateKey(privateKey string) error {
	privateKeyDecoded, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return fmt.Errorf("unable to decode Base64 string: %w", err)
	}

	block, _ := pem.Decode(privateKeyDecoded)

	enc.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("unable to parse Private Key: %w", err)
	}

	return nil
}

// PrivateKey gets the Private Key value.
func (enc *RSAEncryptor) PrivateKey() string {
	privateKeyBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(enc.privateKey),
		},
	)

	return base64.StdEncoding.EncodeToString(privateKeyBytes)
}

// SetPublicKey sets the Public Key value.
func (enc *RSAEncryptor) SetPublicKey(publicKey string) error {
	publicKeyDecoded, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return fmt.Errorf("unable to decode Base64 string: %w", err)
	}

	block, _ := pem.Decode(publicKeyDecoded)

	ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("unable to parse Public Keu: %w", err)
	}

	if key, ok := ifc.(*rsa.PublicKey); ok {
		enc.publicKey = key

		return nil
	}

	return errUnableToDecodeKey
}

// PublicKey gets the Public Key value.
func (enc *RSAEncryptor) PublicKey() (string, error) {
	ASN1, err := x509.MarshalPKIXPublicKey(enc.publicKey)
	if err != nil {
		return "", fmt.Errorf("unable to marshal PPublic Key: %w", err)
	}

	publicKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: ASN1,
	})

	return base64.StdEncoding.EncodeToString(publicKeyBytes), nil
}

// GenerateKeyPair generates a new key pair.
func (enc *RSAEncryptor) GenerateKeyPair() error {
	keyPair, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("unable to ggenerate RSA Key Pair: %w", err)
	}

	enc.privateKey = keyPair
	enc.publicKey = &keyPair.PublicKey

	return nil
}

// Encrypt encrypts data with given public key.
func (enc *RSAEncryptor) Encrypt(text string) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, enc.publicKey, []byte(text), nil)
	if err != nil {
		return "", fmt.Errorf("unable to encrypt text: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// Decrypt decrypts data with given private key.
func (enc *RSAEncryptor) Decrypt(encrypted string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("unable to decode Base64 string: %w", err)
	}

	decryptedBytes, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, enc.privateKey, encryptedBytes, nil)
	if err != nil {
		return "", fmt.Errorf("unable to decrypt text: %w", err)
	}

	return string(decryptedBytes), nil
}
