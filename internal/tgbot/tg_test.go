package tgbot

import (
	"os"
	"testing"
)

func TestLoadTokenFromEnv(t *testing.T) {
	os.Setenv("TG_API_SECRET_KEY", "test_token")
	defer os.Unsetenv("TG_API_SECRET_KEY")

	token, err := loadTokenFromEnv()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if token != "test_token" {
		t.Errorf("expected token to be 'test_token', but got '%s'", token)
	}
}

func TestLoadTokenFromEnvMissingEnvVar(t *testing.T) {
	_, err := loadTokenFromEnv()
	if err == nil {
		t.Error("expected error, but got nil")
	}
}
