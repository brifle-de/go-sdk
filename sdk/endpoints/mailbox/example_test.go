package mailbox_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/mailbox"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func ExampleSearchMyInbox() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	page := float32(1)
	search := mailbox.InboxSearch{
		Page:   &page,
		Filter: &mailbox.InboxSearchFilter{Subject: sdk.String("Invoice")},
	}

	res, respStatus, err := mailbox.SearchMyInbox(client, ctx, &search)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		for _, item := range res.Results {
			fmt.Println(*item.Id, *item.Subject)
		}
	}
}
