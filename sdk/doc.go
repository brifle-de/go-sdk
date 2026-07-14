// Package sdk is the entry point of the Brifle Go SDK: it creates an
// authenticated client for the Brifle API and provides small helpers for
// building request values.
//
// The SDK handles authentication for you. You supply an API key and secret
// once; the client obtains an access token on the first request and renews it
// transparently before it expires — you never handle tokens directly.
//
// # Creating a client
//
//	import (
//		"context"
//		"time"
//
//		"github.com/brifle-de/brifle-sdk/sdk"
//		"github.com/brifle-de/brifle-sdk/sdk/endpoints/status"
//		"github.com/brifle-de/brifle-sdk/sdk/middleware"
//	)
//
//	credentials := middleware.Credentials{
//		ApiKey:    "your-api-key",
//		ApiSecret: "your-api-secret",
//	}
//
//	client, err := sdk.NewClient("https://sandbox-api.brifle.de", credentials)
//	if err != nil {
//		panic(err)
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
//	defer cancel()
//
//	res, respStatus, err := status.GetStatus(client, ctx)
//
// Available servers: https://sandbox-api.brifle.de (sandbox) and
// https://api.brifle.de (production).
//
// # Calling endpoints
//
// Endpoints live in the sub-packages under sdk/endpoints. Each endpoint is a
// package-level function whose first two arguments are the client and a
// context.Context, and which returns (result, *api.ResponseStatus, error):
//
//	result, respStatus, err := content.GetContent(client, ctx, &documentID, &readFlag)
//	if err != nil {
//		// transport, decoding or invalid-input error
//	}
//	if respStatus.HttpStatus != 200 {
//		// a non-2xx HTTP response is reported here, not via err;
//		// respStatus.ErrorCode carries the Brifle error code
//	}
//
// The endpoint packages are:
//
//	status      – service health and feature discovery (no auth)
//	auth        – login and logout (token renewal is automatic)
//	accounts    – basic account information lookup
//	tenants     – list and fetch the tenants you own
//	content     – send and read documents, receiver checks, delivery status,
//	              certificates, paper-mail preview and cover letters
//	mailbox     – search your inbox and outbox
//	signatures  – create signature references and export signatures
//	wallet      – issue, read and revoke wallet items (experimental)
//	address     – parse free-form addresses into structured components
//
// # Pointer helpers
//
// Most request fields are pointers. Use the helpers in this package to build
// them concisely:
//
//	sdk.String("letter")            // *string, nil for ""
//	sdk.Base64Encode(fileBytes)     // *string, base64 of raw bytes
//	sdk.Base64EncodeString("hello") // (*string, error)
//
// For more examples see the docs directory (docs/README.md) in the repository.
package sdk
