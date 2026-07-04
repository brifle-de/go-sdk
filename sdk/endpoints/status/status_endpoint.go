package status

import (
	"context"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// GetStatus retrieves the current status of the Brifle API, including the
// service name, version, availability status and the list of enabled features.
// This endpoint does not require authentication.
func GetStatus(client *sdkClient.BrifleClient, ctx context.Context) (*StatusResponse, *api.ResponseStatus, error) {
	response, err := client.ApiClient.WebApiControllerStatusControllerGetStatus(ctx)
	var res StatusResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// StatusResponse describes the health and capabilities of the Brifle API.
type StatusResponse struct {
	*api.StatusResponse
}
