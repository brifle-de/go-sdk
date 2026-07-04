package tenants_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/tenants"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func ExampleGetMyTenants() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, respStatus, err := tenants.GetMyTenants(client, ctx)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		for _, t := range res.Tenants {
			fmt.Println("tenant:", *t.Id)
		}
	}
}
