# Address

Parse a free-form address string into structured components (street, house number, postcode, city,
country).

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/address`

## ParseAddress

```go
func ParseAddress(client *client.BrifleClient, ctx context.Context, address *string) (*ParsedAddress, *api.ResponseStatus, error)
```

Parses a single address into its best structured interpretation.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

res, respStatus, err := address.ParseAddress(client, ctx, sdk.String("Hauptstraße 5A, 12345 Berlin, Germany"))
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("street: ", res.Street, res.HouseNumber)
fmt.Println("city:   ", res.Postcode, res.City)
fmt.Println("country:", res.Country)
```

## ParseAndExpandAddress

```go
func ParseAndExpandAddress(client *client.BrifleClient, ctx context.Context, address *string) (*ParsedAddressList, *api.ResponseStatus, error)
```

Parses an address and returns **all** plausible structured interpretations (e.g. different spellings
of a street name). Useful when you want to present the user a choice or match against several
variants.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

res, respStatus, err := address.ParseAndExpandAddress(client, ctx, sdk.String("Hauptstrasse 5A, 12345 Berlin"))
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

for _, a := range res.Addresses {
	fmt.Printf("%s %s, %s %s (%s)\n", a.Street, a.HouseNumber, a.Postcode, a.City, a.Country)
}
```
