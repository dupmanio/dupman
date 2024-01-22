package vault

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/otel"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/google/uuid"
	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
	"go.uber.org/zap"
)

type Vault struct {
	logger    *zap.Logger
	client    *vault.Client
	ot        *otel.OTel
	config    Config
	authToken *vault.Secret
}

func New(conf Config, logger *zap.Logger, ot *otel.OTel) (*Vault, error) {
	ctx := context.Background()
	config := vault.DefaultConfig()
	config.Address = conf.Address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault Client: %w", err)
	}

	vaultInst := &Vault{
		logger: logger,
		client: client,
		config: conf,
		ot:     ot,
	}

	if vaultInst.authToken, err = vaultInst.login(ctx); err != nil {
		return nil, err
	}

	return vaultInst, nil
}

func (inst *Vault) login(ctx context.Context) (*vault.Secret, error) {
	appRoleAuth, err := auth.NewAppRoleAuth(
		inst.config.RoleID,
		&auth.SecretID{
			FromString: inst.config.SecretID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize App Role Auth: %w", err)
	}

	authInfo, err := inst.client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate client: %w", err)
	}

	return authInfo, nil
}

func (inst *Vault) CreateUserTransitKey(ctx context.Context, userID uuid.UUID) error {
	path := fmt.Sprintf("transit/keys/user-%s", userID)
	data := map[string]interface{}{
		"type": "aes256-gcm96",
	}

	_, span := inst.ot.GetSpanForFunctionCall(
		ctx,
		1,
		otel.VaultPath(path),
	)
	defer span.End()

	if _, err := inst.client.Logical().Write(path, data); err != nil {
		return fmt.Errorf("unable to create User Transit Key: %w", err)
	}

	return nil
}

func (inst *Vault) EncryptWithUserTransitKey(ctx context.Context, userID uuid.UUID, plaintext string) (string, error) {
	path := fmt.Sprintf("transit/encrypt/user-%s", userID)
	data := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	_, span := inst.ot.GetSpanForFunctionCall(
		ctx,
		1,
		otel.VaultPath(path),
	)
	defer span.End()

	secret, err := inst.client.Logical().Write(path, data)
	if err != nil {
		return "", fmt.Errorf("unable to Encrypt text: %w", err)
	}

	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", domainErrors.ErrSomethingWentWrong
	}

	return ciphertext, nil
}

func (inst *Vault) DecryptWithUserTransitKey(ctx context.Context, userID uuid.UUID, ciphertext string) (string, error) {
	path := fmt.Sprintf("transit/decrypt/user-%s", userID)
	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}

	_, span := inst.ot.GetSpanForFunctionCall(
		ctx,
		1,
		otel.VaultPath(path),
	)
	defer span.End()

	secret, err := inst.client.Logical().Write(path, data)
	if err != nil {
		return "", fmt.Errorf("unable to Decrypt cipher: %w", err)
	}

	plaintext, err := base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
	if err != nil {
		return "", fmt.Errorf("unable to decode plaintext: %w", err)
	}

	return string(plaintext), nil
}
