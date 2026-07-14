package sdk_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/status"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
	"github.com/joho/godotenv"
)

func loadEnv(t *testing.T) {
	// Load the .env.test file from the repository root (one level above the
	// sdk package directory).
	cwd, _ := os.Getwd()
	envFile := filepath.Join(cwd, "../.env.test")
	if err := godotenv.Load(envFile); err != nil {
		t.Skip("No .env.test file found or failed to load")
	}
}

// TestNewClientAuthenticatesAgainstBackend verifies that a client created with
// sdk.NewClient authenticates against the backend automatically. The auth
// middleware renews its token lazily on the first request, so a successful
// authenticated call proves the initial authentication worked.
func TestNewClientAuthenticatesAgainstBackend(t *testing.T) {
	loadEnv(t)

	credentials := middleware.Credentials{
		ApiKey:    os.Getenv("API_KEY"),
		ApiSecret: os.Getenv("API_SECRET"),
	}

	client, err := sdk.NewClient(os.Getenv("ENDPOINT"), credentials)
	if err != nil {
		t.Errorf("Failed to create Brifle client: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Triggers the middleware's lazy authentication against the backend.
	res, respStatus, err := status.GetStatus(client, ctx)
	if err != nil {
		t.Errorf("Authenticated request failed: %v", err)
		return
	}
	if respStatus != nil && respStatus.HttpStatus != 200 {
		t.Errorf("Expected HTTP status 200, got %d", respStatus.HttpStatus)
		return
	}
	if res == nil {
		t.Error("Expected a non-nil response after authentication, got nil")
		return
	}
}

// TestNewClientWithOptsAuthenticatesAgainstBackend covers the same flow through
// the options constructor sdk.NewClientWithOpts.
func TestNewClientWithOptsAuthenticatesAgainstBackend(t *testing.T) {
	loadEnv(t)

	credentials := middleware.Credentials{
		ApiKey:    os.Getenv("API_KEY"),
		ApiSecret: os.Getenv("API_SECRET"),
	}

	client, err := sdk.NewClientWithOpts(os.Getenv("ENDPOINT"), credentials, &sdk.ClientOps{
		SkipTlsVerification: true,
	})
	if err != nil {
		t.Errorf("Failed to create Brifle client: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, respStatus, err := status.GetStatus(client, ctx)
	if err != nil {
		t.Errorf("Authenticated request failed: %v", err)
		return
	}
	if respStatus != nil && respStatus.HttpStatus != 200 {
		t.Errorf("Expected HTTP status 200, got %d", respStatus.HttpStatus)
		return
	}
	if res == nil {
		t.Error("Expected a non-nil response after authentication, got nil")
		return
	}
}
