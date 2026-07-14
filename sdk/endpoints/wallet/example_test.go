package wallet_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/wallet"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func ExampleCreateWalletItem() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
	req := wallet.CreateWalletRequest{
		Type:    wallet.ProofOfOwnership,
		Subject: "Membership Card",
		Data: []wallet.DataElement{
			{
				Name:        sdk.String("Member Name"),
				Value:       sdk.String("Max Mustermann"),
				Type:        sdk.String(wallet.TypeText),
				ReferenceId: sdk.String("member_name"),
			},
		},
	}

	// CreateWalletItem returns HTTP 201 on success.
	res, respStatus, err := wallet.CreateWalletItem(client, ctx, &tenant, &req)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 201 {
		fmt.Println("created wallet item:", *res.Id)
	}
}
