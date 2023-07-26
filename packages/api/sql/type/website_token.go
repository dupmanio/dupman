package sqltype

import (
	"context"
	"errors"
	"fmt"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/encryptor"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	errUnableToEncryptToken = errors.New("unable to encrypt token")
	errInvalidEncryptionKey = errors.New("invalid encryption key")
)

type WebsiteToken string

func (token WebsiteToken) Decrypt(privateKey string) (string, error) {
	rsaEncryptor := encryptor.NewRSAEncryptor()

	err := rsaEncryptor.SetPrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to set Private Key: %w", err)
	}

	decrypted, err := rsaEncryptor.Decrypt(string(token))
	if err != nil {
		return "", fmt.Errorf("unable to decrypt text: %w", err)
	}

	return decrypted, nil
}

func (token WebsiteToken) Encrypt(publicKey string) (string, error) {
	rsaEncryptor := encryptor.NewRSAEncryptor()
	if err := rsaEncryptor.SetPublicKey(publicKey); err != nil {
		return "", fmt.Errorf("unable to set Public Key: %w", err)
	}

	encrypted, err := rsaEncryptor.Encrypt(string(token))
	if err != nil {
		return "", fmt.Errorf("unable to encrypt text: %w", err)
	}

	return encrypted, nil
}

func (token WebsiteToken) GormValue(ctx context.Context, tx *gorm.DB) clause.Expr {
	if encryptionKey, ok := ctx.Value(constant.EncryptionKeyKey).(string); ok {
		if encrypted, err := token.Encrypt(encryptionKey); err == nil {
			return clause.Expr{SQL: "?", Vars: []interface{}{encrypted}}
		}

		_ = tx.AddError(errUnableToEncryptToken)

		return clause.Expr{}
	}

	_ = tx.AddError(errInvalidEncryptionKey)

	return clause.Expr{}
}
