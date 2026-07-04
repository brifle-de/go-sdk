package content

import (
	"context"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// GetDeliveryStatus retrieves the delivery status of a document by its ID. The
// response distinguishes between "brifle" (electronic) and "physical" (paper
// mail) delivery modes.
func GetDeliveryStatus(client *sdkClient.BrifleClient, ctx context.Context, documentId *string) (*DeliveryStatus, *api.ResponseStatus, error) {
	if documentId == nil || *documentId == "" {
		return nil, nil, errors.New("document ID is required")
	}
	response, err := client.ApiClient.WebApiControllerContentControllerGetDeliveryStatus(ctx, *documentId)
	var res DeliveryStatus
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// DeliveryStatus is the delivery status of a document.
type DeliveryStatus struct {
	*api.ContentGetDeliveryStatusResponse
}

// CheckReceiverBulk checks whether multiple receivers exist on Brifle in a
// single request. The order of the results matches the order of the input.
func CheckReceiverBulk(client *sdkClient.BrifleClient, ctx context.Context, receivers *[]ReceiverData) (*ReceiverBulkCheckResponse, *api.ResponseStatus, error) {
	if receivers == nil {
		return nil, nil, errors.New("receivers is nil")
	}

	converted := make([]api.ApiSendContentReceiverRequest, 0, len(*receivers))
	for i := range *receivers {
		r := buildReceiver(&(*receivers)[i])
		if r == nil {
			return nil, nil, errors.New("one or more receivers are invalid")
		}
		converted = append(converted, *r)
	}

	request := api.ApiSendContentReceiverBulkRequest{Receivers: &converted}
	response, err := client.ApiClient.WebApiControllerContentControllerCheckReceiverBulk(ctx, request)
	var res ReceiverBulkCheckResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// ReceiverBulkCheckResponse holds the results of a bulk receiver check.
type ReceiverBulkCheckResponse struct {
	*api.ReceiverBulkExistResponse
}
