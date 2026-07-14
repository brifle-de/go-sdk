# Status

Check the health and capabilities of the Brifle API. This endpoint does **not** require
authentication.

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/status`

## GetStatus

```go
func GetStatus(client *client.BrifleClient, ctx context.Context) (*StatusResponse, *api.ResponseStatus, error)
```

Returns the service name, version, availability status and the list of enabled features.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

res, respStatus, err := status.GetStatus(client, ctx)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("status:  ", *res.Status)
fmt.Println("service: ", *res.Service)
fmt.Println("version: ", *res.Version)
if res.Features != nil {
	fmt.Println("features:", *res.Features)
}
```
