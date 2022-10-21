package data

import (
	"database/sql"
	"strings"

	"github.com/accrescent/devportal/quality"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "devportal.db?_fk=yes&_journal=WAL")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitializeDB(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		gh_id INT NOT NULL,
		access_token TEXT NOT NULL,
		expiry_time INT NOT NULL
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		gh_id INT PRIMARY KEY,
		email TEXT NOT NULL
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS reviewers (
		user_gh_id INT PRIMARY KEY REFERENCES users(gh_id) ON DELETE CASCADE,
		email TEXT NOT NULL
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS usable_email_cache (
		session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
		email TEXT NOT NULL,
		PRIMARY KEY (session_id, email)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_apps (
		id TEXT NOT NULL,
		session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL,
		path TEXT NOT NULL,
		PRIMARY KEY (id, session_id)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS review_errors (
		id TEXT PRIMARY KEY
	) STRICT`); err != nil {
		return err
	}
	if err := populateReviewErrors(db); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_app_review_errors (
		staging_app_id TEXT NOT NULL,
		staging_app_session_id TEXT NOT NULL,
		review_error_id TEXT NOT NULL REFERENCES review_errors(id) ON DELETE CASCADE,
		PRIMARY KEY (staging_app_id, staging_app_session_id, review_error_id),
		FOREIGN KEY (staging_app_id, staging_app_session_id)
			REFERENCES staging_apps(id, session_id)
			ON DELETE CASCADE
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_update_review_errors (
		staging_app_id INT NOT NULL REFERENCES staging_app_updates(id) ON DELETE CASCADE,
		review_error_id TEXT NOT NULL REFERENCES review_errors(id) ON DELETE CASCADE,
		PRIMARY KEY (staging_app_id, review_error_id)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS submitted_apps (
		id TEXT PRIMARY KEY,
		gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL,
		path TEXT NOT NULL
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS submitted_app_review_errors (
		submitted_app_id TEXT NOT NULL REFERENCES submitted_apps(id) ON DELETE CASCADE,
		review_error_id TEXT NOT NULL REFERENCES review_errors(id) ON DELETE CASCADE,
		PRIMARY KEY (submitted_app_id, review_error_id)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS published_apps (
		id TEXT PRIMARY KEY,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS user_permissions (
		app_id TEXT NOT NULL REFERENCES published_apps(id) ON DELETE CASCADE,
		user_gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
		can_update INT NOT NULL CHECK(can_update in (FALSE, TRUE)) DEFAULT FALSE,
		PRIMARY KEY (app_id, user_gh_id)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_app_updates (
		id INTEGER PRIMARY KEY,
		app_id TEXT NOT NULL REFERENCES published_apps(id) ON DELETE CASCADE,
		session_id TEXT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL,
		path TEXT NOT NULL,
		UNIQUE (app_id, version_code)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS submitted_updates (
		id INTEGER PRIMARY KEY,
		app_id TEXT NOT NULL REFERENCES published_apps(id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL,
		path TEXT NOT NULL,
		UNIQUE (app_id, version_code)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS submitted_update_review_errors (
		submitted_app_id INT NOT NULL REFERENCES submitted_updates(id) ON DELETE CASCADE,
		review_error_id TEXT NOT NULL REFERENCES review_errors(id) ON DELETE CASCADE,
		PRIMARY KEY (submitted_app_id, review_error_id)
	) STRICT`); err != nil {
		return err
	}

	return nil
}

func populateReviewErrors(db *sql.DB) error {
	query := "INSERT OR IGNORE INTO review_errors (id) VALUES "
	var inserts []string
	var params []interface{}
	for _, reviewError := range quality.PermissionReviewBlacklist {
		inserts = append(inserts, "(?)")
		params = append(params, reviewError)
	}
	query = query + strings.Join(inserts, ",")
	if _, err := db.Exec(query, params...); err != nil {
		return err
	}

	return nil
}
