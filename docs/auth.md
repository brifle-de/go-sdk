# Authentication

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/auth`

> In normal usage you do **not** need to call these functions. When you create a client with
> `sdk.NewClient`, the SDK logs in automatically and renews the access token before it expires.
> Use these functions only if you need to manage tokens yourself.

## Login

```go
func Login(client *client.BrifleClient, ctx context.Context, apiKey string, apiSecret string) (*LoginResponse, *api.ResponseStatus, error)
```

Exchanges an API key and secret for an access token.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

res, respStatus, err := auth.Login(client, ctx, "your-api-key", "your-api-secret")
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("login failed: http %d", respStatus.HttpStatus)
}

fmt.Println("access token:", *res.AccessToken)
fmt.Println("expires in:  ", *res.ExpiresIn)
```

## Logout

```go
func Logout(client *client.BrifleClient, ctx context.Context, token *string) (*api.ResponseStatus, error)
```

Revokes an access token so it can no longer be used.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

login, _, err := auth.Login(client, ctx, "your-api-key", "your-api-secret")
if err != nil {
	log.Fatal(err)
}

respStatus, err := auth.Logout(client, ctx, login.AccessToken)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("logout failed: http %d", respStatus.HttpStatus)
}
```
