package auth_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/auth"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

// You rarely need Login directly — sdk.NewClient logs in and renews tokens for
// you. This shows the manual flow.
func ExampleLogin() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, respStatus, err := auth.Login(client, ctx, "your-api-key", "your-api-secret")
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		fmt.Println("access token:", *res.AccessToken)
	}
}
