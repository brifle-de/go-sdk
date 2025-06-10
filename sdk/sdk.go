package sdk

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	apiClient "github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/auth"
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

// String returns a pointer to the string value if it is not empty, otherwise returns nil.
func String(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}

// Base64Encode encodes a string to Base64 and returns a pointer to the encoded string.
func Base64Encode(bytes []byte) *string {
	if len(bytes) == 0 {
		return nil
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return &encoded
}
func Base64EncodeString(str string) (*string, error) {
	if str == "" {
		return nil, errors.New("input string cannot be empty")
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return &encoded, nil
}
