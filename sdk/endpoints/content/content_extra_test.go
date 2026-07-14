package content_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/content"
)

func TestGetDeliveryStatus(t *testing.T) {
	brifleClient := getClient(t)
	documentId := os.Getenv("TEST_DOC_ID_CERTIFICATE")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, status, err := content.GetDeliveryStatus(brifleClient, ctx, &documentId)
	if err != nil {
		t.Errorf("GetDeliveryStatus failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if res == nil || res.ContentGetDeliveryStatusResponse == nil {
		t.Error("GetDeliveryStatus response is nil")
		return
	}
}

func TestCheckReceiverBulk(t *testing.T) {
	brifleClient := getClient(t)

	firstName := os.Getenv("TEST_RECEIVER_FIRST_NAME")
	lastName := os.Getenv("TEST_RECEIVER_LAST_NAME")
	dateOfBirth := os.Getenv("TEST_RECEIVER_DATE_OF_BIRTH")
	placeOfBirth := os.Getenv("TEST_RECEIVER_PLACE_OF_BIRTH")

	receivers := []content.ReceiverData{
		{
			BirthInformation: &content.BirthInformationReceiver{
				FirstName:    &firstName,
				LastName:     &lastName,
				DateOfBirth:  &dateOfBirth,
				PlaceOfBirth: &placeOfBirth,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, status, err := content.CheckReceiverBulk(brifleClient, ctx, &receivers)
	if err != nil {
		t.Errorf("CheckReceiverBulk failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if res == nil || res.ReceiverBulkExistResponse == nil {
		t.Error("CheckReceiverBulk response is nil")
		return
	}
}

func TestListCoverLetters(t *testing.T) {
	brifleClient := getClient(t)
	tenant := os.Getenv("TEST_TENANT")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, status, err := content.ListCoverLetters(brifleClient, ctx, &tenant)
	if err != nil {
		t.Errorf("ListCoverLetters failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if res == nil || res.CoverLettersOverviewResponse == nil {
		t.Error("ListCoverLetters response is nil")
		return
	}
}

func TestPreviewPaperMail(t *testing.T) {
	brifleClient := getClient(t)
	tenant := os.Getenv("TEST_TENANT")
	pathToFile := "./test/welcome.pdf"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fileContent, err := os.ReadFile(pathToFile)
	if err != nil {
		t.Errorf("Failed to read file: %v", err)
		return
	}

	req := content.PreviewPaperMailRequest{
		To: &content.PreviewReceiver{
			AddressLine1: sdk.String("Test Street 1"),
			City:         sdk.String("Berlin"),
			Country:      sdk.String("DE"),
			PostalCode:   sdk.String("12345"),
		},
		CoverLetter: &content.PreviewCoverLetter{
			Enable: true,
			Name:   sdk.String("default"),
			Type:   sdk.String(content.CoverLetterDefault),
		},
		Body: &content.PreviewBody{
			Content: sdk.Base64Encode(fileContent),
			Type:    sdk.String("application/pdf"),
		},
	}

	pdf, status, err := content.PreviewPaperMail(brifleClient, ctx, &tenant, &req)
	if err != nil {
		t.Errorf("PreviewPaperMail failed: %v", err)
		return
	}
	if status != nil && status.HttpStatus != 200 {
		t.Errorf("Expected status code 200, got %d", status.HttpStatus)
		return
	}
	if len(pdf) == 0 {
		t.Error("PreviewPaperMail returned empty PDF")
		return
	}
}
