package tenants

import (
	"context"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// Gets a tenant by its ID.
func GetTenant(client *sdkClient.BrifleClient, context context.Context, tenantId *string) (*TenantResponse, *api.ResponseStatus, error) {
	if tenantId == nil {
		return nil, nil, errors.New("tenantId cannot be nil")
	}

	response, err := client.ApiClient.WebApiControllerTenantControllerGetTenant(context, *tenantId)
	var res TenantResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// Gets all tenants that the account owns.
func GetMyTenants(client *sdkClient.BrifleClient, context context.Context) (*MyTenantsResponse, *api.ResponseStatus, error) {
	response, err := client.ApiClient.WebApiControllerTenantControllerGetOwn(context)
	var res MyTenantsResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

type TenantResponse struct {
	*api.TenantResponse
}

type MyTenantsResponse struct {
	Total   int64             `json:"total"`
	Tenants []*TenantResponse `json:"tenants"`
}
