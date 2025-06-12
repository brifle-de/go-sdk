package mailbox_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/mailbox"
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

func TestSearchInbox(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	page := float32(1)
	req := &mailbox.InboxSearch{
		Page: &page,
		Filter: &mailbox.InboxSearchFilter{
			Subject: sdk.String("test"),
			State:   []*string{sdk.String("read")},
			Type:    sdk.String("letter"),
		},
	}

	inboxRes, status, err := mailbox.SearchMyInbox(brifleClient, ctx, req)
	if err != nil {
		t.Errorf("Failed to search inbox: %v", err)
		return
	}

	if status == nil || status.HttpStatus != 200 {
		t.Errorf("Expected status 200, got %d", status)
		return
	}

	if inboxRes != nil {
		t.Error("Expected messages, but got none")
	}

}

func TestOutbox(t *testing.T) {
	brifleClient := getClient(t)
	if brifleClient == nil {
		t.Skip("Skipping test due to client creation failure")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tenants := os.Getenv("TEST_TENANT")

	page := float32(1)
	req := &mailbox.OutboxSearch{
		Page:       &page,
		SenderUser: sdk.String(os.Getenv("TEST_USER_ID")),
		Filter: &mailbox.OutboxFilter{
			Subject: sdk.String("test"),
			State:   []*string{sdk.String("active")},
			Type:    sdk.String("letter"),
		},
	}

	outboxRes, status, err := mailbox.SearchOutbox(brifleClient, ctx, &tenants, req)
	if err != nil {
		t.Errorf("Failed to search outbox: %v", err)
		return
	}

	if status == nil || status.HttpStatus != 200 {
		t.Errorf("Expected status 200, got %d", status)
		return
	}

	if outboxRes == nil {
		t.Error("Expected messages, but got none")
	}
}
