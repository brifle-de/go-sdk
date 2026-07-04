package auth

import (
	"context"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// Login authenticates a user using an API key and secret.
func Login(client *sdkClient.BrifleClient, context context.Context, apiKey string, apiSecret string) (*LoginResponse, *api.ResponseStatus, error) {
	loginRequest := api.LoginApiKeyRequest{
		Key:    &apiKey,
		Secret: &apiSecret,
	}
	response, err := client.ApiClient.WebApiControllerAuthControllerCreate(context, loginRequest)
	var res LoginResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// Logout revokes the given access token so it can no longer be used. Note that
// the SDK manages tokens automatically; use this only if you need to explicitly
// invalidate a token you obtained yourself.
func Logout(client *sdkClient.BrifleClient, ctx context.Context, token *string) (*api.ResponseStatus, error) {
	if token == nil || *token == "" {
		return nil, errors.New("token is required")
	}
	request := api.RevokeTokenRequest{Token: *token}
	response, err := client.ApiClient.WebApiControllerAuthControllerRevoke(ctx, request)
	if err != nil {
		return nil, err
	}
	status, _, err := api.ParseResponseAsString(response)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// Types

type LoginResponse struct {
	*api.LoginResponse
}
