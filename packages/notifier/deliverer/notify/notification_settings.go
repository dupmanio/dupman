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
		"WebsiteNeedsUpdates": {
			Title:   "Your Website is under risk!",
			Type:    "website_needs_updates",
			Message: "Some modules on your website are outdated.",
		},
		"WebsiteScanningFailed": {
			Title:   "Scanning Failed",
			Type:    "website_scanning_failed",
			Message: "We are not able to scan your website. It might be down.",
		},
	}
}
