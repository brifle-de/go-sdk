# Brifle Go SDK — Documentation

The Brifle Go SDK is a typed wrapper around the [Brifle API](https://brifle.de). It handles
authentication and token renewal for you, and converts JSON responses into Go structs.

The SDK is organized into small, versioned packages under `sdk/endpoints/*`. Each endpoint is a
plain package-level function that takes a `*client.BrifleClient` and a `context.Context`, and
returns `(result, *api.ResponseStatus, error)`. This keeps the public surface stable even when the
underlying OpenAPI names change.

## Contents

| Package | Description |
|---|---|
| [Status](status.md) | Service health and feature discovery (no auth). |
| [Authentication](auth.md) | Login and logout. Token renewal is automatic. |
| [Accounts](accounts.md) | Basic account information lookup. |
| [Tenants](tenants.md) | List and fetch the tenants you own. |
| [Content](content.md) | Send documents, read them, check receivers, delivery status/certificates, paper-mail preview. |
| [Cover Letters](cover-letters.md) | Manage cover letter templates for physical delivery. |
| [Mailbox](mailbox.md) | Search your inbox and outbox. |
| [Signatures](signatures.md) | Create signature references and export signatures. |
| [Wallet](wallet.md) | Issue, read and revoke wallet items (experimental). |
| [Address](address.md) | Parse free-form addresses into structured components. |

## Installation

The published module path (`github.com/brifle-de/brifle-sdk`) differs from the GitHub repository
(`github.com/brifle-de/go-sdk`), so add a `replace` directive to your `go.mod`:

```go
replace github.com/brifle-de/brifle-sdk => github.com/brifle-de/go-sdk v0.0.1

require github.com/brifle-de/brifle-sdk v0.0.1
```

Then:

```bash
go get ./...
```

## Servers

| Environment | Base URL |
|---|---|
| Sandbox | `https://sandbox-api.brifle.de` |
| Production | `https://api.brifle.de` |

## Creating a client

You authenticate with an API key and secret. The SDK obtains an access token on the first request
and transparently renews it before it expires — you never handle tokens directly.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/status"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func main() {
	credentials := middleware.Credentials{
		ApiKey:    "your-api-key",
		ApiSecret: "your-api-secret",
	}

	client, err := sdk.NewClient("https://sandbox-api.brifle.de", credentials)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, respStatus, err := status.GetStatus(client, ctx)
	if err != nil {
		panic(err)
	}
	if respStatus.HttpStatus != 200 {
		panic(fmt.Sprintf("unexpected status: %d", respStatus.HttpStatus))
	}
	fmt.Println("service:", *res.Service, "version:", *res.Version)
}
```

### Client options

`NewClientWithOpts` lets you pass options — for example to skip TLS verification against a local
sandbox:

```go
client, err := sdk.NewClientWithOpts(endpoint, credentials, &sdk.ClientOps{
	SkipTlsVerification: true,
})
```

## Return values and error handling

Every endpoint returns three values:

```go
result, respStatus, err := content.GetContent(client, ctx, &documentId, &readFlag)
```

- `err` (`error`) — a transport, encoding or input-validation failure. Always check this first.
- `respStatus` (`*api.ResponseStatus`) — the HTTP result. `respStatus.HttpStatus` is the HTTP status
  code; `respStatus.ErrorCode` is the Brifle error code on non-2xx responses. **A non-2xx HTTP
  response is reported here, not via `err`**, so always check `respStatus.HttpStatus` too.
- `result` — the typed response (nil on failure).

```go
if err != nil {
	// network / decoding / bad input
	return err
}
if respStatus.HttpStatus != 200 {
	return fmt.Errorf("brifle error %d (http %d)", respStatus.ErrorCode, respStatus.HttpStatus)
}
// use result
```

## Helpers

Most request fields are pointers. The `sdk` package provides helpers to build them:

```go
sdk.String("letter")              // *string, returns nil for ""
sdk.Base64Encode(fileBytes)       // *string, base64 of raw bytes
sdk.Base64EncodeString("hello")   // (*string, error)
```

## Contexts

Pass a `context.Context` with a timeout to every call:

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()
```
