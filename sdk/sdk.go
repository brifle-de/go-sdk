package sdk

import (
	"context"
	"net/http"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	apiClient "github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/auth.go"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

func NewClient(server string, credentials middleware.Credentials) (*apiClient.BrifleClient, error) {
	client, err := api.NewClient(server)
	if err != nil {
		return nil, err
	}

	brifle_client := &apiClient.BrifleClient{
		ApiClient: client,
	}

	// add middleware to http client
	client.Client = &http.Client{
		Transport: &middleware.AuthTransport{
			State: middleware.BrifleClientState{
				AuthInterval:      3600, // 1 hour in seconds
				LastAuthenticated: 0,    // initial value
				Token:             "",   // initial value
				Credentials:       &credentials,
			},
			AllowTokenRenewal: true, // allow token renewal
			RenewToken: func() (string, error) {
				res, status, err := auth.Login(brifle_client, context.TODO(), credentials.ApiKey, credentials.ApiSecret)
				if err != nil && status.HttpStatus != http.StatusOK {
					return "", err
				}
				return *res.AccessToken, nil
			},
		},
	}

	return brifle_client, nil
}
