package model

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/encryptor"
	"gorm.io/gorm"
)

type KeyPair struct {
	Base

	PrivateKey string
	PublicKey  string
}

func (entity *KeyPair) BeforeCreate(tx *gorm.DB) error {
	if err := entity.Base.BeforeCreate(tx); err != nil {
		return fmt.Errorf("unable to run parent BeforeCreate Callback: %w", err)
	}

	rsaEncryptor := encryptor.NewRSAEncryptor()

	err := rsaEncryptor.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("unable to generate RSA Key Pair: %w", err)
	}

	entity.PrivateKey = rsaEncryptor.PrivateKey()
	entity.PublicKey, err = rsaEncryptor.PublicKey()

	if err != nil {
		return fmt.Errorf("unable to get RSA Public Key: %w", err)
	}

	return nil
}
