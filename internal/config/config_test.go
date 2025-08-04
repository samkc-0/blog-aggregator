package config

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config := Read()
	if !strings.Contains(config.DbUrl, "postgres") {
		t.Errorf("Expected config.DbUrl to be a postgres url, got %s", config.DbUrl)
	}
}
