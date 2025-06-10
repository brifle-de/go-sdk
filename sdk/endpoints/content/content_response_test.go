package content_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/content"
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

func TestReceiverExists(t *testing.T) {
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
	firstName := os.Getenv("TEST_RECEIVER_FIRST_NAME")
	lastName := os.Getenv("TEST_RECEIVER_LAST_NAME")
	dateOfBirth := os.Getenv("TEST_RECEIVER_DATE_OF_BIRTH")
	placeOfBirth := os.Getenv("TEST_RECEIVER_PLACE_OF_BIRTH")

	receiverData := content.ReceiverData{
		BirthInformation: &content.BirthInformationReceiver{
			FirstName:    &firstName,
			LastName:     &lastName,
			DateOfBirth:  &dateOfBirth,
			PlaceOfBirth: &placeOfBirth,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Test ReceiverExists
	receiverRes, status, err := content.CheckReceiver(brifleClient, ctx, &receiverData)

	if err != nil {
		t.Errorf("ReceiverExists failed: %v", err)
		return
	}
	if status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if *receiverRes.Receiver.Type != "birth_info" {
		t.Errorf("Expected receiver type 'birth_info', got '%s'", *receiverRes.Receiver.Type)
		return
	}
}

func TestSendContent(t *testing.T) {

	loadEnv(t)

	credentials := middleware.Credentials{
		ApiKey:    os.Getenv("API_KEY"),
		ApiSecret: os.Getenv("API_SECRET"),
	}
	brifleClient, err := sdk.NewClient("https://internaltest-api.brifle.de", credentials)
	if err != nil {
		t.Errorf("Failed to create Brifle client: %v", err)
		return
	}
	firstName := os.Getenv("TEST_RECEIVER_FIRST_NAME")
	lastName := os.Getenv("TEST_RECEIVER_LAST_NAME")
	dateOfBirth := os.Getenv("TEST_RECEIVER_DATE_OF_BIRTH")
	placeOfBirth := os.Getenv("TEST_RECEIVER_PLACE_OF_BIRTH")
	pathToFile := "./test/welcome.pdf"
	tenant := os.Getenv("TEST_TENANT")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// read file content
	fileContent, err := os.ReadFile(pathToFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}
	// convert file content to base64
	fileContentBase64 := sdk.Base64Encode(fileContent)
	if err != nil {
		t.Errorf("Failed to encode file content to base64: %v", err)
		return
	}

	document1 := content.ContentItem{
		Content: fileContentBase64,
		Type:    sdk.String("application/pdf"),
	}

	req := content.SendContentRequest{
		To: &content.ReceiverData{
			BirthInformation: &content.BirthInformationReceiver{
				FirstName:    &firstName,
				LastName:     &lastName,
				DateOfBirth:  &dateOfBirth,
				PlaceOfBirth: &placeOfBirth,
			},
		},
		Type:    sdk.String(content.Letter),
		Body:    &[]content.ContentItem{document1},
		Subject: sdk.String("Welcome to Brifle from Go!"),
	}

	res, status, err := content.SendContent(brifleClient, ctx, &tenant, &req)
	if err != nil {
		t.Errorf("SendContent failed: %v", err)
		return
	}
	if status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if res == nil {
		t.Error("SendContent response is nil")
		return
	}

}

func TestGetDocument(t *testing.T) {
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

	documentId := "53C9084932FA27B068424A5FCA81974873E54BC88AAB3B5CCB45C4E6E2C90BB1"
	readFlag := false

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Test authentication
	docRes, status, err := content.GetContent(brifleClient, ctx, &documentId, &readFlag)

	if err != nil {
		t.Errorf("Authentication failed: %v", err)
		return
	}

	if status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status)
		return
	}

	if docRes == nil {
		t.Error("Authentication response is nil")
		return
	}
	// check if content array is empty
	if len(*docRes.Content) == 0 {
		t.Error("Content array is empty")
		return
	}
}
