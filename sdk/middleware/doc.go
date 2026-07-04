// Package middleware provides the HTTP transport that authenticates Brifle API
// requests.
//
// [Credentials] holds your API key and secret. [AuthTransport] is an
// http.RoundTripper that injects a bearer token into every request and renews
// it automatically before it expires. Both are wired up for you by
// [github.com/brifle-de/brifle-sdk/sdk.NewClient]; you only need to supply
// Credentials.
package middleware
