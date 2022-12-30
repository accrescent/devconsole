package data

import (
	"golang.org/x/oauth2"

	"github.com/accrescent/devconsole/config"
)

type DB interface {
	Open(dsn string) error
	Initialize() error
	LoadConfig() (*oauth2.Config, *config.Config, error)
	Close() error

	CreateSession(ghID int64, accessToken string) (id string, err error)
	GetSessionInfo(id string) (ghId int64, accessToken string, err error)
	DeleteExpiredSessions() error
	DeleteSession(id string) error

	CanUserRegister(ghID int64) (bool, error)
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
	GetStagingAppInfo(appID string, ghID int64) (appHandle string, iconHandle string, err error)
	GetSubmittedAppInfo(
		appID string,
	) (
		app App,
		ghID int64,
		iconID int,
		issueGroupID *int,
		appHandle string,
		iconHandle string,
		err error,
	)
	GetSubmittedApps(ghID int64) ([]App, error)
	ApproveApp(appID string) error
	PublishApp(appID string) error
	SubmitApp(appID string, label string, ghID int64) error
	DeleteSubmittedApp(appID string) error

	CreateUpdate(app AppWithIssues, ghID int64, fileHandle string) error
	GetSubmittedUpdates(ghID int64) ([]AppWithIssues, error)
	GetUpdateInfo(
		appID string,
		versionCode int,
	) (firstVersion int, versionName string, fileHandle string, issueGroupID *int, err error)
	GetUpdates(reviewerGhID int64) ([]AppWithIssues, error)
	GetStagingUpdateInfo(
		appID string,
		versionCode int32,
		ghID int64,
	) (
		label string,
		versionName string,
		fileHandle string,
		issueGroupID *int,
		needsReview bool,
		err error,
	)
	ApproveUpdate(appID string, versionCode int, versionName string, issueGroupID *int) error
	SubmitUpdate(app App, fileHandle string, issueGroupID *int, needsReview bool) error
	DeleteSubmittedUpdate(appID string, versionCode int) error
}
