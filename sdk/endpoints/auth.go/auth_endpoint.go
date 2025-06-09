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
	if err != nil {
		return nil, nil, err
	}
	if response == nil {
		return nil, nil, errors.New("login response is nil")
	}
	// parse body into LoginResponse
	if response.Body == nil {
		return nil, nil, errors.New("login response body is nil")
	}
	var res LoginResponse
	status, err := api.ParseResponse(response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// Types

type LoginResponse struct {
	*api.LoginResponse
}
