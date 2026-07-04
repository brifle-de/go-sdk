# Wallet

Issue, read and revoke wallet items. A wallet item is a verifiable credential (proof of ownership or
permission) that can be assigned to a user and optionally exported to Apple/Google Wallet.

> **Experimental.** The wallet API is experimental — contact the Brifle team before relying on it in
> production.

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/wallet`

## Constants

```go
wallet.ProofOfOwnership  // "proof_of_ownership"
wallet.ProofOfPermission // "proof_of_permission"

wallet.TypeText // "text"
wallet.TypeUrl  // "url"
wallet.TypeJwt  // "jwt"
```

## CreateWalletItem

```go
func CreateWalletItem(client *client.BrifleClient, ctx context.Context, tenant *string, req *CreateWalletRequest) (*WalletMeta, *api.ResponseStatus, error)
```

Issues a new wallet item. **On success the API responds with HTTP `201`.**

```go
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
		{
			Name:        sdk.String("Member ID"),
			Value:       sdk.String("A-12345"),
			Type:        sdk.String(wallet.TypeText),
			ReferenceId: sdk.String("member_id"),
		},
	},
	// Optional display styling.
	Styles: &wallet.WalletStyle{
		BackgroundColor: sdk.String("#F0F0F0"),
		PrimaryColor:    sdk.String("#FF5733"),
		Rows: []wallet.Row{
			{Left: "member_name", Right: "member_id"},
		},
	},
	// Optional export to external wallets.
	ExportWallet: &wallet.ExportWallet{
		AppleWallet:  boolPtr(true),
		GoogleWallet: boolPtr(true),
	},
}

res, respStatus, err := wallet.CreateWalletItem(client, ctx, &tenant, &req)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 201 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("created wallet item:", *res.Id)
```

`boolPtr` is a tiny helper (`func boolPtr(b bool) *bool { return &b }`); the SDK does not ship one for
`bool`.

### `CreateWalletRequest` fields

| Field | Type | Required | Description |
|---|---|---|---|
| `Type` | `string` | yes | `wallet.ProofOfOwnership` or `wallet.ProofOfPermission`. |
| `Subject` | `string` | yes | Human-readable subject, e.g. "Driver's License". |
| `Data` | `[]DataElement` | yes | The data elements stored in the item. |
| `Styles` | `*WalletStyle` | no | Display styling; a default is used if omitted. |
| `ExportWallet` | `*ExportWallet` | no | Allow export to Apple/Google Wallet. |
| `Immutable` | `*bool` | no | If true, the item cannot be changed after creation. |
| `RetentionPeriodDays` | `*int` | no | Days to retain an unassigned item (default 30). |
| `ExpiresAt` | `*time.Time` | no | Expiration; nil means no expiration. |
| `NotBefore` | `*time.Time` | no | Item is not valid before this time. |

## ReadWalletItem

```go
func ReadWalletItem(client *client.BrifleClient, ctx context.Context, tenant *string, id *string) (*WalletItem, *api.ResponseStatus, error)
```

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
id := "00205a17-865f-4390-95e9-4b5feb146afc"

res, respStatus, err := wallet.ReadWalletItem(client, ctx, &tenant, &id)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("subject:", *res.Meta.Subject)
if res.Data != nil && res.Data.Elements != nil {
	for _, el := range *res.Data.Elements {
		fmt.Printf("  %s = %s\n", *el.Name, *el.Value)
	}
}
```

## RevokeWalletItem

```go
func RevokeWalletItem(client *client.BrifleClient, ctx context.Context, tenant *string, id *string) (*api.ResponseStatus, error)
```

Revokes a wallet item. An item can only be revoked while it has not yet been assigned to a user
wallet.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
id := "00205a17-865f-4390-95e9-4b5feb146afc"

respStatus, err := wallet.RevokeWalletItem(client, ctx, &tenant, &id)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}
```
