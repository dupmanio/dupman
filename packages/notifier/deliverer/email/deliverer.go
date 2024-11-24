package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/notifier/config"
	"gopkg.in/mail.v2"
)

type Deliverer struct {
	config                      *config.Config
	dialer                      *mail.Dialer
	notificationSettingsMapping NotificationSettingsMapping
}

type TemplateData struct {
	UserContactInfo     *dto.ContactInfo
	NotificationMessage dto.NotificationMessage
	DeliveryType        NotificationSettings
}

func New(config *config.Config) (*Deliverer, error) {
	dialer := mail.NewDialer(config.Mailer.Host, config.Mailer.Port, config.Mailer.Username, config.Mailer.Password)

	// @todo: embed templates. Maybe move it to config.
	// @todo: configure TLS.

	return &Deliverer{
		config:                      config,
		dialer:                      dialer,
		notificationSettingsMapping: getNotificationSettingsMapping(),
	}, nil
}

func (del *Deliverer) Name() string {
	return "EmailDeliverer"
}

func (del *Deliverer) Deliver(message dto.NotificationMessage, contactInfo *dto.ContactInfo) error {
	notificationSettings, ok := del.notificationSettingsMapping[message.Type]
	if !ok {
		return domainErrors.ErrUnsupportedNotificationType
	}

	messageBody, err := del.composeMessage(message.Type, TemplateData{
		UserContactInfo:     contactInfo,
		NotificationMessage: message,
		DeliveryType:        notificationSettings,
	})
	if err != nil {
		return err
	}

	if err = del.sendEmail(
		del.config.Email.From,
		contactInfo.Email,
		notificationSettings.Subject,
		messageBody,
	); err != nil {
		return fmt.Errorf("unable to send email: %w", err)
	}

	return nil
}

func (del *Deliverer) composeMessage(notificationType string, data TemplateData) (string, error) {
	tmpl, err := template.ParseFiles("packages/notifier/deliverer/email/templates/" + notificationType + ".html")
	if err != nil {
		return "", fmt.Errorf("unable to parse template file: %w", err)
	}

	var tpl bytes.Buffer
	if err = tmpl.Execute(&tpl, data); err != nil {
		return "", fmt.Errorf("unable to render template: %w", err)
	}

	return tpl.String(), nil
}

func (del *Deliverer) sendEmail(from, to, subject, message string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", message)

	if err := del.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("unable to Dial or Send: %w", err)
	}

	return nil
}
