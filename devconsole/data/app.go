package data

type App struct {
	AppID       string `json:"app_id"`
	Label       string `json:"label"`
	VersionCode int32  `json:"version_code"`
	VersionName string `json:"version_name"`
}

type AppWithIssues struct {
	App
	Issues []string `json:"issues,omitempty"`
}
