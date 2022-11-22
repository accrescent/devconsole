package data

import "testing"

func TestSQLite_Open(t *testing.T) {
	s := new(SQLite)
	if err := s.Open(":memory:"); err != nil {
		t.Fatal("failed to open database:", err)
	}

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
