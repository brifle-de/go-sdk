package accounts

import (
	"context"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// GetBasicInformation retrieves basic account information for a given account ID.
func GetBasicInformation(client *sdkClient.BrifleClient, context context.Context, accountId *string) (*AccountBasicInfo, *api.ResponseStatus, error) {
	if accountId == nil {
		return nil, nil, errors.New("accountId is required")
	}
	response, err := client.ApiClient.WebApiControllerAccountsControllerGetBasicInfo(context, *accountId)
	var res AccountBasicInfo
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

type AccountBasicInfo struct {
	*api.BasicAccountInfoResponse
}
