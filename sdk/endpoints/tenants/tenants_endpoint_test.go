package tenants_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/tenants"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
	"github.com/joho/godotenv"
)

func getClient(t *testing.T) *client.BrifleClient {
	// Load .env file
	loadEnv(t)

	// get env variables for credentials
	credentials := middleware.Credentials{
		ApiKey:    os.Getenv("API_KEY"),
		ApiSecret: os.Getenv("API_SECRET"),
	}
	brifleClient, err := sdk.NewClient(os.Getenv("ENDPOINT"), credentials)
	if err != nil {
		t.Errorf("Failed to create Brifle client: %v", err)
		return nil
	}
	return brifleClient
}

func loadEnv(t *testing.T) {
	// Load .env file
	cwd, _ := os.Getwd()
	env_file := filepath.Join(cwd, "../../../.env.test")
	if err := godotenv.Load(env_file); err != nil {
		t.Skip("No .env file found or failed to load")
	}
}

func TestGetTenants(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}
	tenantsId := os.Getenv("TEST_TENANT_ID")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tenants, status, err := tenants.GetTenant(brifleClient, ctx, &tenantsId)
	if err != nil {
		t.Errorf("Failed to get tenants: %v", err)
		return
	}

	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected HTTP status 200, got %d", status.HttpStatus)
		return
	}

	if tenants == nil {
		t.Error("Expected tenants to be non-nil, got nil")
		return
	}

}

func TestGetMyTenants(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tenants, status, err := tenants.GetMyTenants(brifleClient, ctx)
	if err != nil {
		t.Errorf("Failed to get my tenants: %v", err)
		return
	}

	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected HTTP status 200, got %d", status.HttpStatus)
		return
	}

	if tenants == nil {
		t.Error("Expected tenants to be non-nil, got nil")
		return
	}
}
