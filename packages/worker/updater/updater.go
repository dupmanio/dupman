package updater

import (
	"fmt"

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

	totalPages := 1
	// @todo: make this process parallel.
	for currentPage := 1; currentPage <= totalPages; currentPage++ {
		websites, pager, err := upd.systemService.GetWebsites(publicKey, currentPage)
		if err != nil {
			return fmt.Errorf("unable to get Websites: %w", err)
		}

		totalPages = pager.TotalPages
		upd.logger.Info(
			"Fetched single page of websites",
			zap.Int("currentPage", currentPage),
			zap.Int("totalPages", totalPages),
		)

		for _, website := range *websites {
			if err = upd.processWebsite(website); err != nil {
				upd.logger.Error("Unable to process website", zap.Error(err))
			} else {
				upd.logger.Info("Website has been processed successfully")
			}
		}
	}

	return nil
}

func (upd *Updater) processWebsite(website dto.WebsiteOnSystemResponse) error {
	upd.logger.Info(
		"Started processing Website",
		zap.String("websiteID", website.ID.String()),
		zap.String("websiteURL", website.URL),
	)

	token, err := upd.encryptor.Decrypt(website.Token)
	if err != nil {
		return fmt.Errorf("unable to decrypt Website token: %w", err)
	}

	updates, err := upd.fetcher.Fetch(website.URL, token)
	if err != nil {
		return fmt.Errorf("unable to fetch Website Updates: %w", err)
	}

	if err = upd.createWebsiteUpdates(website.ID, updates); err != nil {
		return fmt.Errorf("unable to create Website Updates: %w", err)
	}

	return nil
}

func (upd *Updater) createWebsiteUpdates(websiteID uuid.UUID, updatesRaw []model.Update) error {
	var updates dto.Updates

	_ = copier.Copy(&updates, &updatesRaw)

	if _, err := upd.systemService.CreateWebsiteUpdates(websiteID, &updates); err != nil {
		return fmt.Errorf("unable to create Website Updates: %w", err)
	}

	return nil
}
