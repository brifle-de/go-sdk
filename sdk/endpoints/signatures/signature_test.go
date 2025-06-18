package signatures_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/signatures"
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

func TestCreateSignatureReference(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tenantId := os.Getenv("TEST_TENANT_ID")

	signaturesReferenceOptions := signatures.SignatureReferenceOptions{
		Fields: []signatures.SignatureReferenceField{
			{
				Name:    "Test field",
				Purpose: "Test purpose",
				Role:    "Test role",
			},
		},
	}

	signatureRef, status, err := signatures.CreateSignatureReference(brifleClient, ctx, &tenantId, &signaturesReferenceOptions)

	if err != nil {
		t.Errorf("Failed to create signature reference: %v", err)
		return
	}

	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Unexpected HTTP status: %d", status.HttpStatus)
		return
	}

	if signatureRef == nil {
		t.Error("Signature reference is nil")
		return
	}

	if signatureRef.Id == "" {
		t.Error("Signature reference ID is empty")
		return
	}

}

func TestExportSignature(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	signatureId := os.Getenv("EXPORT_SIGNATURE_ID") // Replace with a valid signature ID

	signatureRes, status, err := signatures.ExportSignature(brifleClient, ctx, &signatureId, &signatures.ExportOptions{
		Format: "xml"})

	if err != nil {
		t.Errorf("Failed to export signature: %v", err)
		return
	}

	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Unexpected HTTP status: %d", status.HttpStatus)
		return
	}

	if signatureRes == nil {
		t.Error("Exported signature is nil")
		return
	}

	// make sure signatureRes is at least 1000 characters long
	if len(*signatureRes) < 1000 {
		t.Error("Exported signature is too short")
		return
	}

}
