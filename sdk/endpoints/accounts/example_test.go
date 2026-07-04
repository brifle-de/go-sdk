package accounts_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/accounts"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func ExampleGetBasicInformation() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	accountID := "2802510314782548"
	res, respStatus, err := accounts.GetBasicInformation(client, ctx, &accountID)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 && res.LastName != nil {
		fmt.Println("last name:", *res.LastName)
	}
}
