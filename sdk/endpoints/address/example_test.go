package address_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/address"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func ExampleParseAddress() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, respStatus, err := address.ParseAddress(client, ctx, sdk.String("Hauptstraße 5A, 12345 Berlin, Germany"))
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		fmt.Printf("%s %s, %s %s (%s)\n", res.Street, res.HouseNumber, res.Postcode, res.City, res.Country)
	}
}
