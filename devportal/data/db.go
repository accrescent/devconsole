package data

type DB interface {
	Open() error
	Initialize() error
	DeleteExpiredSessions() error
	CreateSession(id string, ghID int64, accessToken string) error
	GetUserRoles(ghID int64) (registered bool, reviewer bool, err error)
	ApproveApp(appID string) error
	GetUpdateInfo(
		appID string,
		versionCode int,
	) (firstVersion int, versionName string, path string, err error)
	ApproveUpdate(appID string, versionCode int, versionName string) error
	GetApprovedApps() ([]App, error)
	GetApps(ghID int64) ([]App, error)
	GetPendingApps(reviewerGhID int64) ([]App, error)
	GetUpdates(reviewerGhID int64) ([]App, error)
	DeleteSession(id string) error
	CreateApp(
		id string,
		ghID int64,
		label string,
		versionCode int32,
		versionName string,
		appPath string,
		iconPath string,
		iconHash string,
		issues []string,
	) error
	GetAppInfo(appID string) (versionCode int, err error)
	CreateUpdate(
		id string,
		ghID int64,
		label string,
		versionCode int32,
		versionName string,
		path string,
		issues []string,
	) error
	GetSubmittedAppInfo(
		appID string,
	) (
		ghID int64,
		label string,
		versionCode int,
		versionName string,
		iconID int,
		path string,
		err error,
	)
	PublishApp(
		appID string,
		label string,
		versionCode int,
		versionName string,
		iconID int,
		ghID int64,
	) error
	CreateUser(ghID int64, email string) error
	DeleteSubmittedApp(appID string) error
	DeleteSubmittedUpdate(appID string, versionCode int) error
	SubmitApp(appID string, ghID int64) error
	GetSessionInfo(id string) (ghId int64, accessToken string, err error)
	GetUserPermissions(appID string, ghID int64) (update bool, err error)
	GetStagingUpdateInfo(
		appID string,
		versionCode int,
		ghID int64,
	) (label string, versionName string, path string, issueGroupID *int, err error)
	SubmitUpdate(
		appID string,
		label string,
		versionCode int,
		versionName string,
		path string,
		issueGroupID *int,
	) error
}
