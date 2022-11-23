package data

import "testing"

func testOpenSQLite(t testing.TB) *SQLite {
	s := new(SQLite)
	if err := s.Open(":memory:"); err != nil {
		t.Fatal("failed to open database:", err)
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
	s := testOpenSQLite(t)
	defer s.Close()

	if err := s.Initialize(); err != nil {
		t.Fatal("failed to initialize database:", err)
	}
}

func TestSQLiteClose(t *testing.T) {
	s := testOpenSQLite(t)
	defer s.Close()
	if err := s.Initialize(); err != nil {
		t.Fatal("failed to initialize database:", err)
	}

	if err := s.CreateUser(123456, "example@example.com"); err != nil {
		t.Fatal("failed to create user")
	}

	s.Close()

	if _, _, err := s.GetUserRoles(123456); err == nil {
		t.Error("query after close succeeded")
	}
}
