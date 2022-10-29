package data

import "database/sql"

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
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_apps (
		id TEXT NOT NULL,
		user_gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL,
		path TEXT NOT NULL,
		issue_group_id INT REFERENCES issue_groups(id),
		PRIMARY KEY (id, user_gh_id)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS issue_groups (
		id INTEGER PRIMARY KEY
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS issues (
		id TEXT NOT NULL,
		issue_group_id INT NOT NULL REFERENCES issue_groups(id),
		PRIMARY KEY (id, issue_group_id)
	) STRICT`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS submitted_apps (
		id TEXT PRIMARY KEY,
		gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL,
		issue_group_id INT REFERENCES issue_groups(id),
		reviewer_gh_id INT NOT NULL REFERENCES reviewers(user_gh_id),
		approved INT NOT NULL CHECK(approved in (FALSE, TRUE)) DEFAULT FALSE,
		path TEXT NOT NULL
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
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS staging_updates (
		id INTEGER PRIMARY KEY,
		app_id TEXT NOT NULL REFERENCES published_apps(id) ON DELETE CASCADE,
		user_gh_id INT NOT NULL REFERENCES users(gh_id) ON DELETE CASCADE,
		label TEXT NOT NULL,
		version_code INT NOT NULL,
		version_name TEXT NOT NULL,
		path TEXT NOT NULL,
		issue_group_id INT REFERENCES issue_groups(id),
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
		reviewer_gh_id INT NOT NULL REFERENCES reviewers(user_gh_id),
		path TEXT NOT NULL,
		issue_group_id INT NOT NULL REFERENCES issue_groups(id),
		UNIQUE (app_id, version_code)
	) STRICT`); err != nil {
		return err
	}

	return nil
}
