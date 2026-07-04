# Tenants

A tenant represents a sending identity (a person or organization) you can send content on behalf of.
Many endpoints require a tenant ID.

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/tenants`

## GetMyTenants

```go
func GetMyTenants(client *client.BrifleClient, ctx context.Context) (*MyTenantsResponse, *api.ResponseStatus, error)
```

Lists all tenants the authenticated account owns.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

res, respStatus, err := tenants.GetMyTenants(client, ctx)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("total tenants:", res.Total)
for _, tnt := range res.Tenants {
	fmt.Println("tenant:", *tnt.Id)
}
```

## GetTenant

```go
func GetTenant(client *client.BrifleClient, ctx context.Context, tenantId *string) (*TenantResponse, *api.ResponseStatus, error)
```

Fetches a single tenant by its ID.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenantId := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
res, respStatus, err := tenants.GetTenant(client, ctx, &tenantId)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("tenant id:", *res.Id)
```
