package scheduler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/encryptor"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/broker"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/service/system"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger        *zap.Logger
	systemService *system.System
	encryptor     encryptor.Encryptor
	broker        *broker.RabbitMQ
}

func New(conf *config.Config, logger *zap.Logger, broker *broker.RabbitMQ) (*Scheduler, error) {
	cred, err := credentials.NewClientCredentials(conf.Dupman.ClientID, conf.Dupman.ClientSecret, []string{})
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

	return &Scheduler{
		logger:        logger,
		encryptor:     rsaEncryptor,
		systemService: system.New(sess),
		broker:        broker,
	}, nil
}

func (scheduler *Scheduler) Process() error {
	scheduler.logger.Info("Starting Scheduler Process")

	publicKey, err := scheduler.encryptor.PublicKey()
	if err != nil {
		return fmt.Errorf("unable to get public key: %w", err)
	}

	currentPage := 1
	totalPages := 1

	var wg sync.WaitGroup

	for currentPage <= totalPages {
		websites, pager, err := scheduler.systemService.GetWebsites(publicKey, currentPage)
		if err != nil {
			scheduler.logger.Error("unable to get Websites", zap.Error(err))
		}

		totalPages = pager.TotalPages
		scheduler.logger.Info(
			"Fetched single page of websites",
			zap.Int("currentPage", currentPage),
			zap.Int("totalPages", totalPages),
		)
		currentPage++

		for _, website := range *websites {
			wg.Add(1)

			go func(website dto.WebsiteOnSystemResponse) {
				defer wg.Done()

				if err = scheduler.scheduleWebsiteScanning(website); err != nil {
					scheduler.logger.Error(
						"Unable to schedule website scanning",
						zap.Error(err),
						zap.String("websiteID", website.ID.String()),
						zap.String("websiteURL", website.URL),
					)
				} else {
					scheduler.logger.Info(
						"Website scanning has been scheduled successfully",
						zap.String("websiteID", website.ID.String()),
						zap.String("websiteURL", website.URL),
					)
				}
			}(website)
		}
	}
	wg.Wait()

	return nil
}

func (scheduler *Scheduler) scheduleWebsiteScanning(website dto.WebsiteOnSystemResponse) error {
	scheduler.logger.Info(
		"Started processing Website",
		zap.String("websiteID", website.ID.String()),
		zap.String("websiteURL", website.URL),
	)

	token, err := scheduler.encryptor.Decrypt(website.Token)
	if err != nil {
		return fmt.Errorf("unable to decrypt Website token: %w", err)
	}

	msg, err := scheduler.composeMessage(website, token)
	if err != nil {
		return fmt.Errorf("unable to compose message: %w", err)
	}

	err = scheduler.broker.PublishToScanner(msg)
	if err != nil {
		return fmt.Errorf("unable to publish message: %w", err)
	}

	return nil
}

func (scheduler *Scheduler) composeMessage(website dto.WebsiteOnSystemResponse, token string) ([]byte, error) {
	message := dto.ScanWebsiteMessage{
		WebsiteID:    website.ID,
		WebsiteURL:   website.URL,
		WebsiteToken: token,
	}

	return json.Marshal(message)
}
