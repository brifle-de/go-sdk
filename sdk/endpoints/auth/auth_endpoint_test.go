package auth_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/auth"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
	"github.com/joho/godotenv"
)

func loadEnv(t *testing.T) {
	// Load .env file
	cwd, _ := os.Getwd()
	env_file := filepath.Join(cwd, "../../../.env.test")
	if err := godotenv.Load(env_file); err != nil {
		t.Skip("No .env file found or failed to load")
	}
}

func TestAuth(t *testing.T) {
	// Load .env file
	loadEnv(t)

	// get env variables for credentials

	credentials := middleware.Credentials{
		ApiKey:    os.Getenv("API_KEY"),
		ApiSecret: os.Getenv("API_SECRET"),
	}
	brifleClient, err := sdk.NewClient("https://internaltest-api.brifle.de", credentials)
	if err != nil {
		t.Errorf("Failed to create Brifle client: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Test authentication
	loginRes, status, err := auth.Login(brifleClient, ctx, credentials.ApiKey, credentials.ApiSecret)

	if err != nil {
		t.Errorf("Authentication failed: %v", err)
		return
	}

	if status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status)
		return
	}

	if loginRes == nil {
		t.Error("Authentication response is nil")
		return
	}
	if *loginRes.AccessToken == "" {
		t.Error("Authentication token is empty")
		return
	}

}
