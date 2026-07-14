package sdk_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/status"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

// Create a client and make a first call. The client authenticates and renews
// its token automatically, so it can be reused for the lifetime of the program.
func Example() {
	credentials := middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	}

	client, err := sdk.NewClient("https://sandbox-api.brifle.de", credentials)
	if err != nil {
		log.Fatal(err)
	}

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
