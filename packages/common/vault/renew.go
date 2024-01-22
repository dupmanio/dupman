package vault

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/otel"
	vault "github.com/hashicorp/vault/api"
	"go.uber.org/zap"
)

type renewResult uint8

const (
	renewError renewResult = 1 << iota
	exitRequested
	expiringAuthToken // will be revoked soon
)

func (inst *Vault) PeriodicallyRenewLeases(ctx context.Context) error {
	inst.logger.Debug("Starting Vault secrets renew loop")
	defer inst.logger.Debug("Vault secrets renew loop ended")

	for {
		renewed, err := inst.renewLeases(ctx, inst.authToken)
		if err != nil {
			return fmt.Errorf("unable to renew leases: %w", err)
		}

		if renewed&exitRequested != 0 {
			return nil
		}

		if renewed&expiringAuthToken != 0 {
			inst.logger.Debug("Auth token can no longer be renewed")

			if inst.authToken, err = inst.login(ctx); err != nil {
				return fmt.Errorf("login authentication error: %w", err)
			}
		}
	}
}

func (inst *Vault) renewLeases(ctx context.Context, authToken *vault.Secret) (renewResult, error) {
	inst.logger.Debug("Starting Vault leases renew loop")
	defer inst.logger.Debug("Vault secrets leases loop ended")

	authTokenWatcherInput := vault.LifetimeWatcherInput{
		Secret: authToken,
	}

	authTokenWatcher, err := inst.client.NewLifetimeWatcher(&authTokenWatcherInput)
	if err != nil {
		return renewError, fmt.Errorf("unable to initialize auth token lifetime watcher: %w", err)
	}

	go authTokenWatcher.Start()
	defer authTokenWatcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil

		case err = <-authTokenWatcher.DoneCh():
			return expiringAuthToken, err

		case info := <-authTokenWatcher.RenewCh():
			inst.logger.Debug(
				"Auth token has been renewed",
				zap.Int(string(otel.DurationKey), info.Secret.Auth.LeaseDuration),
				zap.String(string(otel.RenewedAtKey), info.RenewedAt.String()),
			)
		}
	}
}
