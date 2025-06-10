package auth

import (
	"context"

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

// Types

type LoginResponse struct {
	*api.LoginResponse
}
