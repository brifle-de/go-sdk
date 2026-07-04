package content_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/content"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

// exampleClient builds a client for the runnable examples below.
func exampleClient() *sdkClient.BrifleClient {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})
	return client
}

// Send a PDF to a recipient identified by their birth information (the most
// precise addressing, recommended for sensitive content).
func ExampleSendContent() {
	client := exampleClient()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
	pdf, err := os.ReadFile("welcome.pdf")
	if err != nil {
		log.Fatal(err)
	}

	req := content.SendContentRequest{
		To: &content.ReceiverData{
			BirthInformation: &content.BirthInformationReceiver{
				FirstName:    sdk.String("Max"),
				LastName:     sdk.String("Mustermann"),
				DateOfBirth:  sdk.String("1999-12-12"),
				PlaceOfBirth: sdk.String("Berlin"),
			},
		},
		Type:    sdk.String(content.Letter),
		Subject: sdk.String("Welcome to Brifle"),
		Body: &[]content.ContentItem{
			{Content: sdk.Base64Encode(pdf), Type: sdk.String("application/pdf")},
		},
	}

	res, respStatus, err := content.SendContent(client, ctx, &tenant, &req)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		fmt.Println("document id:", *res.Id)
	}
}

// Address a recipient by email instead of birth information.
func ExampleSendContent_email() {
	client := exampleClient()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
	req := content.SendContentRequest{
		To: &content.ReceiverData{
			Email: &content.EmailReceiver{
				Email: sdk.String("max@example.com"),
				Name:  sdk.String("Max Mustermann"),
			},
		},
		Type:    sdk.String(content.Letter),
		Subject: sdk.String("Welcome to Brifle"),
		Body: &[]content.ContentItem{
			{Content: sdk.String("SGVsbG8="), Type: sdk.String("application/pdf")},
		},
	}

	_, respStatus, err := content.SendContent(client, ctx, &tenant, &req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("status:", respStatus.HttpStatus)
}

// Check whether a receiver exists on Brifle before sending.
func ExampleCheckReceiver() {
	client := exampleClient()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	receiver := content.ReceiverData{
		BirthInformation: &content.BirthInformationReceiver{
			FirstName:    sdk.String("Max"),
			LastName:     sdk.String("Mustermann"),
			DateOfBirth:  sdk.String("1999-12-12"),
			PlaceOfBirth: sdk.String("Berlin"),
		},
	}

	res, respStatus, err := content.CheckReceiver(client, ctx, &receiver)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		fmt.Println("receiver type:", *res.Receiver.Type)
	}
}

// Fetch a document by ID and mark it as read.
func ExampleGetContent() {
	client := exampleClient()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	documentID := "53C9084932FA27B068424A5FCA81974873E54BC88AAB3B5CCB45C4E6E2C90BB1"
	readFlag := true

	res, respStatus, err := content.GetContent(client, ctx, &documentID, &readFlag)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		for _, part := range *res.Content {
			fmt.Println("content-type:", *part.ContentType)
		}
	}
}
