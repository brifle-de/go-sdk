# Accounts

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/accounts`

## GetBasicInformation

```go
func GetBasicInformation(client *client.BrifleClient, ctx context.Context, accountId *string) (*AccountBasicInfo, *api.ResponseStatus, error)
```

Retrieves basic public information about an account by its ID.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

accountId := "2802510314782548"
res, respStatus, err := accounts.GetBasicInformation(client, ctx, &accountId)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

// Fields depend on account type (private vs. business).
if res.FirstName != nil {
	fmt.Println("first name:", *res.FirstName)
}
if res.LastName != nil {
	fmt.Println("last name: ", *res.LastName)
}
if res.CompanyName != nil {
	fmt.Println("company:   ", *res.CompanyName)
}
```
