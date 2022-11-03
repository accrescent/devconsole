package data

type App struct {
	AppID       string   `json:"app_id"`
	Label       string   `json:"label"`
	VersionCode int      `json:"version_code"`
	VersionName string   `json:"version_name"`
	Issues      []string `json:"issues,omitempty"`
}
