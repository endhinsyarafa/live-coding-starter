package database

import (
	"strings"
	"testing"

	"live-coding-starter/backend/internal/config"
)

func TestBuildDSNOmitsEmptyPassword(t *testing.T) {
	dsn := buildDSN(config.DatabaseConfig{
		Host:    "localhost",
		Port:    "5434",
		User:    "macends",
		Name:    "live_coding_db",
		SSLMode: "disable",
	})

	if strings.Contains(dsn, "password=") {
		t.Fatalf("expected empty password to be omitted, got %q", dsn)
	}

	if !strings.Contains(dsn, "dbname=live_coding_db") {
		t.Fatalf("expected dbname in dsn, got %q", dsn)
	}
}
