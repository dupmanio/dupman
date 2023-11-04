package updater

import (
	"fmt"
	"sync"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/encryptor"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/service/system"
	"github.com/dupmanio/dupman/packages/worker/config"
	"github.com/dupmanio/dupman/packages/worker/fetcher"
	"github.com/dupmanio/dupman/packages/worker/model"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type Updater struct {
	logger        *zap.Logger
	systemService *system.System
	encryptor     encryptor.Encryptor
	fetcher       *fetcher.Fetcher
}

func New(conf *config.Config, logger *zap.Logger) (*Updater, error) {
	cred, err := credentials.NewClientCredentials(conf.Dupman.ClientID, conf.Dupman.ClientSecret, conf.Dupman.Scopes)
	if err != nil {
		return nil, fmt.Errorf("unable to initiate credentials provider: %w", err)
	}

	sess, err := session.New(&dupman.Config{Credentials: cred})
	if err != nil {
		return nil, fmt.Errorf("unable to create dupman session: %w", err)
	}

	rsaEncryptor := encryptor.NewRSAEncryptor()
	if err = rsaEncryptor.GenerateKeyPair(); err != nil {
		return nil, fmt.Errorf("unable to generate RSA Key Pair: %w", err)
	}

	return &Updater{
		logger:        logger,
		encryptor:     rsaEncryptor,
		fetcher:       fetcher.New(),
		systemService: system.New(sess),
	}, nil
}

func (upd *Updater) Process() error {
	upd.logger.Info("Starting Updater Process")

	publicKey, err := upd.encryptor.PublicKey()
	if err != nil {
		return fmt.Errorf("unable to get public key: %w", err)
	}

	currentPage := 1
	totalPages := 1

	var wg sync.WaitGroup

	for currentPage <= totalPages {
		websites, pager, err := upd.systemService.GetWebsites(publicKey, currentPage)
		if err != nil {
			upd.logger.Error("unable to get Websites", zap.Error(err))
		}

		totalPages = pager.TotalPages
		upd.logger.Info(
			"Fetched single page of websites",
			zap.Int("currentPage", currentPage),
			zap.Int("totalPages", totalPages),
		)
		currentPage++

		for _, website := range *websites {
			wg.Add(1)

			go func(website dto.WebsiteOnSystemResponse) {
				defer wg.Done()

				if status, err := upd.processWebsite(website); err != nil {
					upd.logger.Error(
						"Unable to process website",
						zap.Error(err),
						zap.String("websiteID", website.ID.String()),
						zap.String("websiteURL", website.URL),
					)
				} else {
					upd.logger.Info(
						"Website has been processed successfully",
						zap.String("websiteID", website.ID.String()),
						zap.String("websiteURL", website.URL),
						zap.String("statusID", status.ID.String()),
						zap.String("statusInfo", status.Info),
						zap.String("statusState", status.State),
					)
				}
			}(website)
		}
	}
	wg.Wait()

	return nil
}

func (upd *Updater) processWebsite(website dto.WebsiteOnSystemResponse) (*dto.StatusOnSystemResponse, error) {
	upd.logger.Info(
		"Started processing Website",
		zap.String("websiteID", website.ID.String()),
		zap.String("websiteURL", website.URL),
	)

	status := dto.Status{
		State: dto.StatusStateUpToDated,
	}

	token, err := upd.encryptor.Decrypt(website.Token)
	if err != nil {
		upd.logger.Error(
			"unable to decrypt Website token",
			zap.String("websiteID", website.ID.String()),
			zap.String("websiteURL", website.URL),
			zap.Error(err),
		)

		status.State = dto.StatusStateScanningFailed
		status.Info = err.Error()
	}

	updates, err := upd.fetcher.Fetch(website.URL, token)
	if err != nil {
		upd.logger.Error(
			"unable to fetch Website Updates",
			zap.String("websiteID", website.ID.String()),
			zap.String("websiteURL", website.URL),
			zap.Error(err),
		)

		status.State = dto.StatusStateScanningFailed
		status.Info = err.Error()
	}

	if len(updates) != 0 {
		status.State = dto.StatusStateNeedsUpdate
	}

	websiteStatus, err := upd.updateWebsiteStatus(website.ID, status, updates)
	if err != nil {
		return nil, fmt.Errorf("unable to create Website Updates: %w", err)
	}

	return websiteStatus, nil
}

func (upd *Updater) updateWebsiteStatus(
	websiteID uuid.UUID, status dto.Status,
	updatesRaw []model.Update,
) (*dto.StatusOnSystemResponse, error) {
	var updates dto.Updates

	_ = copier.Copy(&updates, &updatesRaw)

	updateResponse, err := upd.systemService.UpdateWebsiteStatus(websiteID, &status, &updates)
	if err != nil {
		return nil, fmt.Errorf("unable to create Website Updates: %w", err)
	}

	return &updateResponse.Status, nil
}
