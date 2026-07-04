package wallet_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/wallet"
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

// TestWalletLifecycle creates a wallet item, reads it back, then revokes it.
func TestWalletLifecycle(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}
	tenant := os.Getenv("TEST_TENANT")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req := wallet.CreateWalletRequest{
		Type:    wallet.ProofOfOwnership,
		Subject: "SDK Test Item",
		Data: []wallet.DataElement{
			{
				Name:        sdk.String("Full Name"),
				Value:       sdk.String("Max Mustermann"),
				Type:        sdk.String(wallet.TypeText),
				ReferenceId: sdk.String("full_name"),
			},
		},
	}

	created, status, err := wallet.CreateWalletItem(brifleClient, ctx, &tenant, &req)
	if err != nil {
		t.Errorf("CreateWalletItem failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 201 {
		t.Errorf("Expected status code 201, got %d", status.HttpStatus)
		return
	}
	if created == nil || created.WalletMetaResponse == nil || created.Id == nil {
		t.Error("CreateWalletItem response is missing an id")
		return
	}

	id := *created.Id

	read, status, err := wallet.ReadWalletItem(brifleClient, ctx, &tenant, &id)
	if err != nil {
		t.Errorf("ReadWalletItem failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if read == nil || read.WalletItemResponse == nil {
		t.Error("ReadWalletItem response is nil")
		return
	}

	status, err = wallet.RevokeWalletItem(brifleClient, ctx, &tenant, &id)
	if err != nil {
		t.Errorf("RevokeWalletItem failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
}
