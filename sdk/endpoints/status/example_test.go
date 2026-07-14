package status_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/status"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func ExampleGetStatus() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, respStatus, err := status.GetStatus(client, ctx)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		fmt.Println("service:", *res.Service, "version:", *res.Version)
	}
}
