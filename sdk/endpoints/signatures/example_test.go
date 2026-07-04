package signatures_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/signatures"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func ExampleCreateSignatureReference() {
	client, _ := sdk.NewClient("https://sandbox-api.brifle.de", middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tenantID := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
	opts := signatures.SignatureReferenceOptions{
		Fields: []signatures.SignatureReferenceField{
			{Name: "signature_1", Purpose: "approval", Role: "customer"},
		},
	}

	res, respStatus, err := signatures.CreateSignatureReference(client, ctx, &tenantID, &opts)
	if err != nil {
		log.Fatal(err)
	}
	if respStatus.HttpStatus == 200 {
		fmt.Println("reference id:", res.Id)
	}
}
