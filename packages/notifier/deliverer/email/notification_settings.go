package email

type NotificationSettings struct {
	Subject string
}

type NotificationSettingsMapping map[string]NotificationSettings

func getNotificationSettingsMapping() NotificationSettingsMapping {
	// @todo: implement better mechanism to configure notifications.
	// @todo: implement feature to allow placeholders in texts.
	return NotificationSettingsMapping{
		"hello": {
			Subject: "Hello, World!",
		},
	}
}
