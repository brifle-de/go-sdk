package address

import (
	"context"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// ParseAddress parses a free-form address string into its structured components
// (street, house number, postcode, city, country).
func ParseAddress(client *sdkClient.BrifleClient, ctx context.Context, address *string) (*ParsedAddress, *api.ResponseStatus, error) {
	if address == nil || *address == "" {
		return nil, nil, errors.New("address is required")
	}
	request := api.ParseAddressRequest{Address: *address}
	response, err := client.ApiClient.WebApiControllerAddressControllerParseAddress(ctx, request)
	var res ParsedAddress
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// ParseAndExpandAddress parses a free-form address string and returns all
// plausible structured interpretations (e.g. different spellings of a street).
func ParseAndExpandAddress(client *sdkClient.BrifleClient, ctx context.Context, address *string) (*ParsedAddressList, *api.ResponseStatus, error) {
	if address == nil || *address == "" {
		return nil, nil, errors.New("address is required")
	}
	request := api.ParseAddressRequest{Address: *address}
	response, err := client.ApiClient.WebApiControllerAddressControllerParseAndExpandAddress(ctx, request)
	var res api.ParseAddressArrayResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	list := ParsedAddressList{Addresses: res}
	return &list, status, nil
}

// ParsedAddress is a single structured address.
type ParsedAddress struct {
	*api.ParseAddressResponse
}

// ParsedAddressList holds multiple structured address interpretations.
type ParsedAddressList struct {
	Addresses api.ParseAddressArrayResponse `json:"addresses"`
}
