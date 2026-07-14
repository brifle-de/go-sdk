package address_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/address"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
	"github.com/joho/godotenv"
)

func getClient(t *testing.T) *client.BrifleClient {
	loadEnv(t)
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
	cwd, _ := os.Getwd()
	env_file := filepath.Join(cwd, "../../../.env.test")
	if err := godotenv.Load(env_file); err != nil {
		t.Skip("No .env file found or failed to load")
	}
}

func TestParseAddress(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	input := sdk.String("Hauptstraße 5A, 12345 Berlin, Germany")
	res, status, err := address.ParseAddress(brifleClient, ctx, input)
	if err != nil {
		t.Errorf("ParseAddress failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if res == nil || res.ParseAddressResponse == nil {
		t.Error("ParseAddress response is nil")
		return
	}
}

func TestParseAndExpandAddress(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	input := sdk.String("Hauptstraße 5A, 12345 Berlin, Germany")
	res, status, err := address.ParseAndExpandAddress(brifleClient, ctx, input)
	if err != nil {
		t.Errorf("ParseAndExpandAddress failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if res == nil {
		t.Error("ParseAndExpandAddress response is nil")
		return
	}
}
