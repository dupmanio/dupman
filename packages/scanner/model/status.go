package model

type CoreCompatibility struct {
	Compatibility string `json:"compatibility"`
	Compatible    string `json:"compatible"`
	Message       string `json:"message"`
}

type Release struct {
	Name              string            `json:"name"`
	Version           string            `json:"version"`
	Tag               string            `json:"tag"`
	Status            string            `json:"status"`
	Link              string            `json:"link"`
	Date              string            `json:"date"`
	Security          string            `json:"security"`
	CoreCompatibility CoreCompatibility `json:"core_compatibility"`
}

type Update struct {
	Name               string    `json:"name"`
	Title              string    `json:"title"`
	Link               string    `json:"link"`
	Type               string    `json:"type"`
	CurrentVersion     string    `json:"current_version"`
	LatestVersion      string    `json:"latest_version"`
	RecommendedVersion string    `json:"recommended_version"`
	Releases           []Release `json:"releases"`
	InstallType        string    `json:"install_type"`
	Status             int       `json:"status"`
}

type Status struct {
	Updates []Update `json:"updates"`
}
