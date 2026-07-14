// Package client defines [BrifleClient], the handle passed to every endpoint
// function in the SDK.
//
// You normally do not construct a BrifleClient directly — use
// [github.com/brifle-de/brifle-sdk/sdk.NewClient] or
// [github.com/brifle-de/brifle-sdk/sdk.NewClientWithOpts], which wire up
// automatic authentication and token renewal.
package client
