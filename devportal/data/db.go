package data

type DB interface {
	Open(dsn string) error
	Initialize() error
	Close() error

	CreateSession(ghID int64, accessToken string) (id string, err error)
	GetSessionInfo(id string) (ghId int64, accessToken string, err error)
	DeleteExpiredSessions() error
	DeleteSession(id string) error

	CreateUser(ghID int64, email string) error
	GetUserPermissions(appID string, ghID int64) (update bool, err error)
	GetUserRoles(ghID int64) (registered bool, reviewer bool, err error)

	CreateReviewer(ghID int64, email string) error

	CreateApp(
		app AppWithIssues,
		ghID int64,
		appFileHandle string,
		iconFileHandle string,
		iconHash string,
	) error
	GetAppInfo(appID string) (versionCode int32, err error)
	GetApprovedApps() ([]App, error)
	GetApps(ghID int64) ([]App, error)
	GetPendingApps(reviewerGhID int64) ([]AppWithIssues, error)
	GetSubmittedAppInfo(appID string) (app App, ghID int64, iconID int, fileHandle string, err error)
	ApproveApp(appID string) error
	PublishApp(appID string) error
	SubmitApp(appID string, label string, ghID int64) error
	DeleteSubmittedApp(appID string) error

	CreateUpdate(app AppWithIssues, ghID int64, fileHandle string) error
	GetUpdateInfo(
		appID string,
		versionCode int,
	) (firstVersion int, versionName string, fileHandle string, err error)
	GetUpdates(reviewerGhID int64) ([]AppWithIssues, error)
	GetStagingUpdateInfo(
		appID string,
		versionCode int,
		ghID int64,
	) (label string, versionName string, fileHandle string, issueGroupID *int, err error)
	ApproveUpdate(appID string, versionCode int, versionName string) error
	SubmitUpdate(app App, fileHandle string, issueGroupID *int) error
	DeleteSubmittedUpdate(appID string, versionCode int) error
}
