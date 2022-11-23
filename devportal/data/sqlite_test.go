package data

import (
	"encoding/hex"
	"testing"
)

func testOpenSQLite(t testing.TB) *SQLite {
	s := new(SQLite)
	if err := s.Open(":memory:"); err != nil {
		t.Fatal("failed to open database:", err)
	}

	return s
}

func testCreateSQLite(t testing.TB) *SQLite {
	s := testOpenSQLite(t)
	if err := s.Initialize(); err != nil {
		t.Fatal("failed to initialize database:", err)
	}

	return s
}

func TestSQLiteOpen(t *testing.T) {
	s := testOpenSQLite(t)
	defer s.Close()

	t.Run("trusted_schema value", func(t *testing.T) {
		var trustedSchema bool
		if err := s.db.QueryRow("PRAGMA trusted_schema").Scan(&trustedSchema); err != nil {
			t.Fatal("failed to read trusted_schema:", err)
		}

		if trustedSchema {
			t.Error("trusted_schema is ON")
		}
	})
}

func TestSQLiteInitialize(t *testing.T) {
	s := testCreateSQLite(t)
	defer s.Close()
}

func TestSQLiteClose(t *testing.T) {
	s := testCreateSQLite(t)
	defer s.Close()

	if err := s.CreateUser(123456, "example@example.com"); err != nil {
		t.Fatal("failed to create user")
	}

	s.Close()

	if _, _, err := s.GetUserRoles(123456); err == nil {
		t.Error("query after close succeeded")
	}
}

func TestSQLiteSession(t *testing.T) {
	s := testCreateSQLite(t)
	defer s.Close()

	var testGHID int64 = 123456
	testToken := "token-1234"
	testSIDLen := 16

	sessionID, err := s.CreateSession(testGHID, testToken)
	if err != nil {
		t.Fatal("failed to create session:", err)
	}

	t.Run("session ID properties", func(t *testing.T) {
		decoded, err := hex.DecodeString(sessionID)
		if err != nil {
			t.Error("session ID is not hex encoded:", err)
		}
		decodedLen := len(decoded)
		if decodedLen != testSIDLen {
			t.Errorf("session ID length is %d but expected %d", decodedLen, testSIDLen)
		}
	})

	t.Run("get", func(t *testing.T) {
		ghID, token, err := s.GetSessionInfo(sessionID)
		if err != nil {
			t.Fatal("failed to get session:", err)
		}

		if ghID != testGHID {
			t.Errorf("GitHub ID is %d but expected %d", ghID, testGHID)
		}
		if token != testToken {
			t.Errorf("access token is %s but expected %s", token, testToken)
		}
	})
}
