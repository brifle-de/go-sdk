package sdk

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	apiClient "github.com/brifle-de/brifle-sdk/sdk/client"
	"github.com/brifle-de/brifle-sdk/sdk/endpoints/auth"
	"github.com/brifle-de/brifle-sdk/sdk/middleware"
)

type ClientOps struct {
	SkipTlsVerification bool // skip TLS verification for the client
}

func NewClientWithOpts(server string, credentials middleware.Credentials, opts *ClientOps) (*apiClient.BrifleClient, error) {
	if opts == nil {
		opts = &ClientOps{SkipTlsVerification: false} // default value
	}
	return newClient(server, credentials, opts)
}

func NewClient(server string, credentials middleware.Credentials) (*apiClient.BrifleClient, error) {
	return newClient(server, credentials, nil)
}

func newClient(server string, credentials middleware.Credentials, opts *ClientOps) (*apiClient.BrifleClient, error) {
	client, err := api.NewClient(server)
	var skipTlsVerification bool
	if opts != nil {
		skipTlsVerification = opts.SkipTlsVerification
	}

	if err != nil {
		return nil, err
	}

	brifle_client := &apiClient.BrifleClient{
		ApiClient: client,
	}

	// add middleware to http client
	client.Client = &http.Client{
		Transport: &middleware.AuthTransport{
			BaseTransport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: skipTlsVerification,
				},
			},
			State: middleware.BrifleClientState{
				AuthInterval:      3600, // 1 hour in seconds
				LastAuthenticated: 0,    // initial value
				Token:             "",   // initial value
				Credentials:       &credentials,
			},
			AllowTokenRenewal: true, // allow token renewal
			RenewToken: func() (string, error) {
				res, status, err := auth.Login(brifle_client, context.TODO(), credentials.ApiKey, credentials.ApiSecret)
				if err != nil {
					if status != nil && status.HttpStatus != http.StatusOK {
						return "", errors.New("failed to retrieve access token")
					}
					return "", err
				}

				if res == nil || res.AccessToken == nil {
					return "", errors.New("failed to retrieve access token")
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
