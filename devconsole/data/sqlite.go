package data

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/accrescent/devconsole/config"
)

func init() {
	sql.Register(
		"sqlite3_hardened",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				_, err := conn.Exec("PRAGMA trusted_schema = OFF", nil)
				return err
			},
		},
	)
}

type SQLite struct {
	db *sql.DB
}

func (s *SQLite) Open(dsn string) error {
	conn, err := sql.Open("sqlite3_hardened", dsn)
	if err != nil {
		return err
	}
	s.db = conn

	return nil
}

func (s *SQLite) Initialize() error {
	var userVersion int
	if err := s.db.QueryRow("PRAGMA user_version").Scan(&userVersion); err != nil {
		return err
	}

	if userVersion < 1 {
		if _, err := s.db.Exec(`CREATE TABLE config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE sessions (
			id TEXT PRIMARY KEY,
			gh_id INT NOT NULL,
			access_token TEXT NOT NULL,
			expiry_time INT NOT NULL
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE users (
			gh_id INT PRIMARY KEY,
			email TEXT NOT NULL
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE reviewers (
			user_gh_id INT PRIMARY KEY REFERENCES users(gh_id) ON DELETE CASCADE,
			email TEXT NOT NULL
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE staging_apps (
			id TEXT NOT NULL,
			user_gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
			label TEXT NOT NULL,
			version_code INT NOT NULL,
			version_name TEXT NOT NULL,
			file_handle TEXT NOT NULL,
			icon_id INT NOT NULL REFERENCES icons(id),
			issue_group_id INT REFERENCES issue_groups(id),
			PRIMARY KEY (id, user_gh_id)
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE icons (
			id INTEGER PRIMARY KEY,
			file_handle TEXT NOT NULL,
			hash TEXT NOT NULL
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE issue_groups (
			id INTEGER PRIMARY KEY
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE issues (
			id TEXT NOT NULL,
			issue_group_id INT NOT NULL REFERENCES issue_groups(id) ON DELETE CASCADE,
			PRIMARY KEY (id, issue_group_id)
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE submitted_apps (
			id TEXT PRIMARY KEY,
			gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
			label TEXT NOT NULL,
			version_code INT NOT NULL,
			version_name TEXT NOT NULL,
			icon_id INT NOT NULL REFERENCES icons(id),
			issue_group_id INT REFERENCES issue_groups(id),
			reviewer_gh_id INT NOT NULL REFERENCES reviewers(user_gh_id),
			approved INT NOT NULL CHECK(approved in (FALSE, TRUE)) DEFAULT FALSE,
			file_handle TEXT NOT NULL
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE published_apps (
			id TEXT PRIMARY KEY,
			label TEXT NOT NULL,
			version_code INT NOT NULL,
			version_name TEXT NOT NULL,
			icon_id INT NOT NULL REFERENCES icons(id),
			issue_group_id INT REFERENCES issue_groups(id)
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE user_permissions (
			app_id TEXT NOT NULL REFERENCES published_apps(id) ON DELETE CASCADE,
			user_gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
			can_update INT NOT NULL CHECK(can_update in (FALSE, TRUE)) DEFAULT FALSE,
			PRIMARY KEY (app_id, user_gh_id)
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE staging_updates (
			app_id TEXT NOT NULL REFERENCES published_apps(id) ON DELETE CASCADE,
			user_gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
			label TEXT NOT NULL,
			version_code INT NOT NULL,
			version_name TEXT NOT NULL,
			file_handle TEXT NOT NULL,
			issue_group_id INT REFERENCES issue_groups(id),
			PRIMARY KEY (app_id, version_code)
		) STRICT`); err != nil {
			return err
		}
		if _, err := s.db.Exec(`CREATE TABLE submitted_updates (
			app_id TEXT NOT NULL REFERENCES published_apps(id) ON DELETE CASCADE,
			label TEXT NOT NULL,
			version_code INT NOT NULL,
			version_name TEXT NOT NULL,
			reviewer_gh_id INT NOT NULL REFERENCES reviewers(user_gh_id),
			file_handle TEXT NOT NULL,
			issue_group_id INT NOT NULL REFERENCES issue_groups(id),
			PRIMARY KEY (app_id, version_code)
		) STRICT`); err != nil {
			return err
		}

		if _, err := s.db.Exec("PRAGMA user_version = 1"); err != nil {
			return err
		}
	}
	if userVersion < 2 {
		if _, err := s.db.Exec(`CREATE TABLE allowed_users (
			gh_id INT PRIMARY KEY
		) STRICT`); err != nil {
			return err
		}

		if _, err := s.db.Exec("PRAGMA user_version = 2"); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLite) LoadConfig() (*oauth2.Config, *config.Config, error) {
	oauth := new(oauth2.Config)
	conf := new(config.Config)
	if err := s.db.QueryRow(`SELECT
		(SELECT value FROM config WHERE key = 'gh_client_id'),
		(SELECT value FROM config WHERE key = 'gh_client_secret'),
		(SELECT value FROM config WHERE key = 'oauth2_redirect_url'),
		(SELECT value FROM config WHERE key = 'signer_gh_id'),
		(SELECT value FROM config WHERE key = 'repo_url'),
		(SELECT value FROM config WHERE key = 'api_key')
	`).Scan(
		&oauth.ClientID,
		&oauth.ClientSecret,
		&oauth.RedirectURL,
		&conf.SignerGitHubID,
		&conf.RepoURL,
		&conf.APIKey,
	); err != nil {
		return nil, nil, err
	}
	oauth.Endpoint = endpoints.GitHub
	oauth.Scopes = []string{"user:email"}

	return oauth, conf, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) CreateSession(ghID int64, accessToken string) (id string, err error) {
	randBytes := make([]byte, 16)
	if _, err = rand.Read(randBytes); err != nil {
		return
	}
	id = hex.EncodeToString(randBytes)

	_, err = s.db.Exec(
		"INSERT INTO sessions (id, gh_id, access_token, expiry_time) VALUES (?, ?, ?, ?)",
		id,
		ghID,
		accessToken,
		time.Now().Add(24*time.Hour).Unix(),
	)

	return
}

func (s *SQLite) GetSessionInfo(id string) (ghID int64, accessToken string, err error) {
	err = s.db.QueryRow(
		"SELECT gh_id, access_token FROM sessions WHERE id = ? AND expiry_time > ?",
		id,
		time.Now().Unix(),
	).Scan(&ghID, &accessToken)

	return
}

func (s *SQLite) DeleteExpiredSessions() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expiry_time < ?", time.Now().Unix())

	return err
}

func (s *SQLite) DeleteSession(id string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", id)

	return err
}

func (s *SQLite) CanUserRegister(ghID int64) (canRegister bool, err error) {
	err = s.db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM allowed_users WHERE gh_id = ?)",
		ghID,
	).Scan(&canRegister)

	return
}

func (s *SQLite) CreateUser(ghID int64, email string) error {
	_, err := s.db.Exec("INSERT INTO users (gh_id, email) VALUES (?, ?)", ghID, email)

	return err
}

func (s *SQLite) GetUserPermissions(appID string, ghID int64) (update bool, err error) {
	err = s.db.QueryRow(
		`SELECT can_update FROM user_permissions
		WHERE app_id = ? AND user_gh_id = ?`,
		appID,
		ghID,
	).Scan(&update)

	return
}

func (s *SQLite) GetUserRoles(ghID int64) (registered bool, reviewer bool, err error) {
	err = s.db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM users WHERE gh_id = ?),
			EXISTS (SELECT 1 FROM reviewers WHERE user_gh_id = ?)`,
		ghID,
		ghID,
	).Scan(&registered, &reviewer)

	return
}

func (s *SQLite) CreateReviewer(ghID int64, email string) error {
	_, err := s.db.Exec("INSERT INTO reviewers (user_gh_id, email) VALUES (?, ?)", ghID, email)

	return err
}

func (s *SQLite) CreateApp(
	app AppWithIssues,
	ghID int64,
	appFileHandle string,
	iconFileHandle string,
	iconHash string,
) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	var issueGroupID *int64
	if len(app.Issues) > 0 {
		res, err := tx.Exec("INSERT INTO issue_groups DEFAULT VALUES")
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		groupID, err := res.LastInsertId()
		issueGroupID = &groupID
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		insertQuery := "INSERT INTO issues (id, issue_group_id) VALUES "
		var inserts []string
		var params []interface{}
		for _, issue := range app.Issues {
			inserts = append(inserts, "(?, ?)")
			params = append(params, issue, groupID)
		}
		insertQuery = insertQuery + strings.Join(inserts, ",")
		if _, err := tx.Exec(insertQuery, params...); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	res, err := tx.Exec(
		"INSERT INTO icons (file_handle, hash) VALUES (?, ?)",
		iconFileHandle,
		iconHash,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	iconID, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec(
		`REPLACE INTO staging_apps (
			id,
			user_gh_id,
			label,
			version_code,
			version_name,
			file_handle,
			icon_id,
			issue_group_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		app.AppID,
		ghID,
		app.Label,
		app.VersionCode,
		app.VersionName,
		appFileHandle,
		iconID,
		issueGroupID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLite) GetAppInfo(appID string) (versionCode int32, err error) {
	err = s.db.QueryRow(
		"SELECT version_code FROM published_apps WHERE id = ?",
		appID,
	).Scan(&versionCode)

	return
}

func (s *SQLite) GetApprovedApps() ([]App, error) {
	rows, err := s.db.Query(
		`SELECT id, label, version_code, version_name
		FROM submitted_apps
		WHERE approved = TRUE`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []App
	for rows.Next() {
		var appID, label, versionName string
		var versionCode int32
		if err := rows.Scan(&appID, &label, &versionCode, &versionName); err != nil {
			return nil, err
		}

		app := App{appID, label, versionCode, versionName}
		apps = append(apps, app)
	}

	return apps, nil
}

func (s *SQLite) GetApps(ghID int64) ([]App, error) {
	rows, err := s.db.Query(
		`SELECT id, label, version_code, version_name
		FROM published_apps
		JOIN user_permissions
		ON user_permissions.app_id = published_apps.id
		WHERE user_permissions.user_gh_id = ?`,
		ghID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []App
	for rows.Next() {
		var appID, label, versionName string
		var versionCode int32
		if err := rows.Scan(&appID, &label, &versionCode, &versionName); err != nil {
			return nil, err
		}

		app := App{appID, label, versionCode, versionName}
		apps = append(apps, app)
	}

	return apps, err
}

func (s *SQLite) GetPendingApps(reviewerGhID int64) ([]AppWithIssues, error) {
	dbApps, err := s.db.Query(
		`SELECT id, label, version_code, version_name, issue_group_id
		FROM submitted_apps
		WHERE reviewer_gh_id = ? AND approved = FALSE`,
		reviewerGhID,
	)
	if err != nil {
		return nil, err
	}
	defer dbApps.Close()
	var apps []AppWithIssues
	for dbApps.Next() {
		var appID, label, versionName string
		var versionCode int32
		var issueGroupID *int
		if err := dbApps.Scan(
			&appID,
			&label,
			&versionCode,
			&versionName,
			&issueGroupID,
		); err != nil {
			return nil, err
		}

		dbIssues, err := s.db.Query(
			"SELECT id FROM issues WHERE issue_group_id = ?",
			issueGroupID,
		)
		if err != nil {
			return nil, err
		}
		defer dbIssues.Close()
		var issues []string
		for dbIssues.Next() {
			var issue string
			if err := dbIssues.Scan(&issue); err != nil {
				return nil, err
			}

			issues = append(issues, issue)
		}

		app := AppWithIssues{App{appID, label, versionCode, versionName}, issues}
		apps = append(apps, app)
	}

	return apps, nil
}

func (s *SQLite) GetStagingAppInfo(appID string, ghID int64) (fileHandle string, err error) {
	err = s.db.QueryRow(
		`SELECT file_handle FROM staging_apps WHERE id = ? AND user_gh_id = ?`,
		appID,
		ghID,
	).Scan(&fileHandle)

	return
}

func (s *SQLite) GetSubmittedAppInfo(
	appID string,
) (
	app App,
	ghID int64,
	iconID int,
	issueGroupID *int,
	fileHandle string,
	err error,
) {
	app.AppID = appID
	err = s.db.QueryRow(
		`SELECT gh_id, label, version_code, version_name, icon_id, issue_group_id, file_handle
		FROM submitted_apps
		WHERE id = ?`,
		appID,
	).Scan(&ghID, &app.Label, &app.VersionCode, &app.VersionName, &iconID, &issueGroupID, &fileHandle)

	return
}

func (s *SQLite) GetSubmittedApps(ghID int64) ([]App, error) {
	rows, err := s.db.Query(
		`SELECT id, label, version_code, version_name
		FROM submitted_apps
		WHERE gh_id = ?`,
		ghID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []App
	for rows.Next() {
		var appID, label, versionName string
		var versionCode int32
		if err := rows.Scan(&appID, &label, &versionCode, &versionName); err != nil {
			return nil, err
		}

		app := App{appID, label, versionCode, versionName}
		apps = append(apps, app)
	}

	return apps, err
}

func (s *SQLite) ApproveApp(appID string) error {
	_, err := s.db.Exec("UPDATE submitted_apps SET approved = TRUE WHERE id = ?", appID)

	return err
}

func (s *SQLite) PublishApp(appID string) error {
	app, ghID, iconID, issueGroupID, _, err := s.GetSubmittedAppInfo(appID)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(
		`INSERT INTO published_apps (
			id,
			label,
			version_code,
			version_name,
			icon_id,
			issue_group_id
		)
		VALUES (?, ?, ?, ?, ?, ?)`,
		app.AppID,
		app.Label,
		app.VersionCode,
		app.VersionName,
		iconID,
		issueGroupID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec(
		"INSERT INTO user_permissions (app_id, user_gh_id, can_update) VALUES (?, ?, TRUE)",
		app.AppID,
		ghID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec("DELETE FROM submitted_apps WHERE id = ?", app.AppID); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLite) SubmitApp(appID string, label string, ghID int64) error {
	var fileHandle, versionName string
	var versionCode, iconID int
	var issueGroupID *int
	if err := s.db.QueryRow(
		`SELECT version_code, version_name, file_handle, icon_id, issue_group_id
		FROM staging_apps
		WHERE id = ? AND user_gh_id = ?`,
		appID,
		ghID,
	).Scan(&versionCode, &versionName, &fileHandle, &iconID, &issueGroupID); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(
		`INSERT INTO submitted_apps (
			id,
			gh_id,
			label,
			version_code,
			version_name,
			icon_id,
			reviewer_gh_id,
			file_handle,
			issue_group_id
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			(SELECT user_gh_id FROM reviewers ORDER BY RANDOM() LIMIT 1),
			?,
			?
		)`,
		appID,
		ghID,
		label,
		versionCode,
		versionName,
		iconID,
		fileHandle,
		issueGroupID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec(
		"DELETE FROM staging_apps WHERE id = ? AND user_gh_id = ?",
		appID,
		ghID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLite) DeleteSubmittedApp(appID string) error {
	_, err := s.db.Exec("DELETE FROM submitted_apps WHERE id = ?", appID)

	return err
}

func (s *SQLite) CreateUpdate(app AppWithIssues, ghID int64, fileHandle string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	var issueGroupID *int64
	if len(app.Issues) > 0 {
		res, err := tx.Exec("INSERT INTO issue_groups DEFAULT VALUES")
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		groupID, err := res.LastInsertId()
		issueGroupID = &groupID
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		insertQuery := "INSERT INTO issues (id, issue_group_id) VALUES "
		var inserts []string
		var params []interface{}
		for _, issue := range app.Issues {
			inserts = append(inserts, "(?, ?)")
			params = append(params, issue, issueGroupID)
		}
		insertQuery = insertQuery + strings.Join(inserts, ",")
		if _, err := tx.Exec(insertQuery, params...); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	if _, err := tx.Exec(
		`REPLACE INTO staging_updates (
			app_id,
			user_gh_id,
			label,
			version_code,
			version_name,
			file_handle,
			issue_group_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		app.AppID,
		ghID,
		app.Label,
		app.VersionCode,
		app.VersionName,
		fileHandle,
		issueGroupID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLite) GetSubmittedUpdates(ghID int64) ([]AppWithIssues, error) {
	dbApps, err := s.db.Query(
		`SELECT
			submitted_updates.app_id,
			submitted_updates.label,
			submitted_updates.version_code,
			submitted_updates.version_name,
			submitted_updates.issue_group_id
		FROM submitted_updates
		JOIN published_apps ON published_apps.id = submitted_updates.app_id
		JOIN user_permissions ON user_permissions.app_id = published_apps.id
		WHERE user_permissions.user_gh_id = ?`,
		ghID,
	)
	if err != nil {
		return nil, err
	}
	defer dbApps.Close()
	var apps []AppWithIssues
	for dbApps.Next() {
		var appID, label, versionName string
		var versionCode int32
		var issueGroupID *int
		if err := dbApps.Scan(
			&appID,
			&label,
			&versionCode,
			&versionName,
			&issueGroupID,
		); err != nil {
			return nil, err
		}

		dbIssues, err := s.db.Query(
			`SELECT issues.id FROM issues WHERE issues.issue_group_id = ?
			AND issues.id NOT IN (
				SELECT issues.id
				FROM published_apps
				JOIN issues
				ON issues.issue_group_id = published_apps.issue_group_id
				AND published_apps.id = ?
			)`,
			issueGroupID,
			appID,
		)
		if err != nil {
			return nil, err
		}
		defer dbIssues.Close()
		var issues []string
		for dbIssues.Next() {
			var issue string
			if err := dbIssues.Scan(&issue); err != nil {
				return nil, err
			}

			issues = append(issues, issue)
		}

		app := AppWithIssues{App{appID, label, versionCode, versionName}, issues}
		apps = append(apps, app)
	}

	return apps, nil
}

func (s *SQLite) GetUpdateInfo(
	appID string,
	versionCode int,
) (firstVersion int, versionName string, fileHandle string, issueGroupID *int, err error) {
	err = s.db.QueryRow(
		`SELECT
			(SELECT MIN(version_code) FROM submitted_updates),
			version_name,
			file_handle,
			issue_group_id
		FROM submitted_updates
		WHERE app_id = ? AND version_code = ?`,
		appID,
		versionCode,
	).Scan(&firstVersion, &versionName, &fileHandle, &issueGroupID)

	return
}

func (s *SQLite) GetUpdates(reviewerGhID int64) ([]AppWithIssues, error) {
	dbApps, err := s.db.Query(
		`SELECT app_id, label, version_code, version_name, issue_group_id
		FROM submitted_updates
		WHERE reviewer_gh_id = ?`,
		reviewerGhID,
	)
	if err != nil {
		return nil, err
	}
	defer dbApps.Close()
	var apps []AppWithIssues
	for dbApps.Next() {
		var appID, label, versionName string
		var versionCode int32
		var issueGroupID *int
		if err := dbApps.Scan(
			&appID,
			&label,
			&versionCode,
			&versionName,
			&issueGroupID,
		); err != nil {
			return nil, err
		}

		dbIssues, err := s.db.Query(
			`SELECT issues.id FROM issues WHERE issues.issue_group_id = ?
			AND issues.id NOT IN (
				SELECT issues.id
				FROM published_apps
				JOIN issues
				ON issues.issue_group_id = published_apps.issue_group_id
				AND published_apps.id = ?
			)`,
			issueGroupID,
			appID,
		)
		if err != nil {
			return nil, err
		}
		defer dbIssues.Close()
		var issues []string
		for dbIssues.Next() {
			var issue string
			if err := dbIssues.Scan(&issue); err != nil {
				return nil, err
			}

			issues = append(issues, issue)
		}

		app := AppWithIssues{App{appID, label, versionCode, versionName}, issues}
		apps = append(apps, app)
	}

	return apps, nil
}

func (s *SQLite) GetStagingUpdateInfo(
	appID string,
	versionCode int32,
	ghID int64,
) (label string, versionName string, fileHandle string, issueGroupID *int, needsReview bool, err error) {
	// We determine the update needs review if it has any review issues that the published app
	// doesn't have, i.e., if the update adds any new review issues.
	err = s.db.QueryRow(
		`SELECT label, version_name, file_handle, issue_group_id, EXISTS(
			SELECT issues.id
			FROM staging_updates
			JOIN issues
			ON issues.issue_group_id = staging_updates.issue_group_id
			WHERE issues.id NOT IN (
				SELECT issues.id
				FROM published_apps
				JOIN issues
				ON issues.issue_group_id = published_apps.issue_group_id
				AND published_apps.id = ?
			)
		)
		FROM staging_updates
		WHERE app_id = ? AND version_code = ? AND user_gh_id = ?`,
		appID,
		appID,
		versionCode,
		ghID,
	).Scan(&label, &versionName, &fileHandle, &issueGroupID, &needsReview)

	return
}

func (s *SQLite) ApproveUpdate(
	appID string,
	versionCode int,
	versionName string,
	issueGroupID *int,
) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(
		`UPDATE published_apps
		SET version_code = ?, version_name = ?, issue_group_id = ?
		WHERE id = ?`,
		versionCode,
		versionName,
		issueGroupID,
		appID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec(
		"DELETE FROM submitted_updates WHERE app_id = ? AND version_code = ?",
		appID,
		versionCode,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLite) SubmitUpdate(app App, fileHandle string, issueGroupID *int, needsReview bool) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if needsReview {
		if _, err := tx.Exec(
			`INSERT INTO submitted_updates (
				app_id,
				label,
				version_code,
				version_name,
				reviewer_gh_id,
				file_handle,
				issue_group_id
			)
			VALUES (
				?,
				?,
				?,
				?,
				(SELECT user_gh_id FROM reviewers ORDER BY RANDOM() LIMIT 1),
				?,
				?
			)`,
			app.AppID,
			app.Label,
			app.VersionCode,
			app.VersionName,
			fileHandle,
			issueGroupID,
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	} else {
		// No review necessary, so publish the update immediately.
		if _, err := tx.Exec(
			`UPDATE published_apps
			SET version_code = ?, version_name = ?, issue_group_id = ?
			WHERE id = ?`,
			app.VersionCode,
			app.VersionName,
			issueGroupID,
			app.AppID,
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	if _, err := tx.Exec(
		"DELETE FROM staging_updates WHERE app_id = ? AND version_code = ?",
		app.AppID,
		app.VersionCode,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SQLite) DeleteSubmittedUpdate(appID string, versionCode int) error {
	_, err := s.db.Exec(
		"DELETE FROM submitted_updates WHERE app_id = ? AND version_code = ?",
		appID,
		versionCode,
	)

	return err
}
