package vault

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/user-api/model"
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

func (inst *Vault) CreateUserTransitKey(ctx context.Context, user *model.User) error {
	path := fmt.Sprintf("transit/keys/user-%s", user.ID)
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
