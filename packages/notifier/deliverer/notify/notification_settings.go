package notify

type NotificationSettings struct {
	Title   string
	Type    string
	Message string
}

type NotificationSettingsMapping map[string]NotificationSettings

func getNotificationSettingsMapping() NotificationSettingsMapping {
	// @todo: implement better mechanism to configure notifications.
	// @todo: implement feature to allow placeholders in texts.
	return NotificationSettingsMapping{
		"hello": {
			Title:   "Hello, World!",
			Type:    "hello",
			Message: "Hello, World!",
		},
	}
}
