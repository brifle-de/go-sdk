// Package auth authenticates against the Brifle API.
//
// In normal usage you do not need this package: a client created with
// [github.com/brifle-de/brifle-sdk/sdk.NewClient] logs in automatically and
// renews its access token before it expires. Use these functions only if you
// want to manage tokens yourself.
//
//	res, respStatus, err := auth.Login(client, ctx, "your-api-key", "your-api-secret")
//	if err == nil && respStatus.HttpStatus == 200 {
//		fmt.Println("access token:", *res.AccessToken)
//	}
//
// See docs/auth.md for more examples.
package auth
