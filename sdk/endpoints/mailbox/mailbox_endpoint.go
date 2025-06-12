package mailbox

import (
	"context"
	"errors"
	"fmt"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// SearchMyInbox searches the inbox of the logged in user
func SearchMyInbox(client *sdkClient.BrifleClient, context context.Context, inboxSearch *InboxSearch) (*InboxSearchResponse, *api.ResponseStatus, error) {
	if inboxSearch == nil {
		return nil, nil, errors.New("inboxSearch cannot be nil")
	}
	if client == nil || client.ApiClient == nil {
		return nil, nil, errors.New("client or client.ApiClient cannot be nil")
	}
	request := inboxSearch.ToMyMailboxRequest()
	response, err := client.ApiClient.WebApiControllerMailboxControllerGetMyInbox(context, *request)
	var res InboxSearchResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	if status != nil && status.HttpStatus != 200 {
		return nil, status, fmt.Errorf("unexpected HTTP status: %d", status.HttpStatus)
	}
	return &res, status, nil
}

// Search the outbox of the account of the logged in user. To get documents sent by a specific user, use the SenderUser field in the OutboxSearch struct.
func SearchOutbox(client *sdkClient.BrifleClient, context context.Context, tenant *string, outboxSearch *OutboxSearch) (*OutboxSearchResponse, *api.ResponseStatus, error) {
	if outboxSearch == nil {
		return nil, nil, errors.New("outboxSearch cannot be nil")
	}
	if client == nil || client.ApiClient == nil {
		return nil, nil, errors.New("client or client.ApiClient cannot be nil")
	}
	request := outboxSearch.ToOutboxRequest(tenant)
	response, err := client.ApiClient.WebApiControllerMailboxControllerGetMyOutbox(context, *tenant, *request)
	var res OutboxSearchResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	if status != nil && status.HttpStatus != 200 {
		return nil, status, fmt.Errorf("unexpected HTTP status: %d", status.HttpStatus)
	}
	return &res, status, nil
}

type InboxSearchFilter struct {
	Subject *string   `json:"subject,omitempty"`
	State   []*string `json:"state,omitempty"`
	Type    *string   `json:"type,omitempty"`
}

type InboxSearch struct {
	// Filter The filter to apply to the search
	Filter *InboxSearchFilter `json:"filter,omitempty"`

	// Page The page number to get
	Page *float32 `json:"page,omitempty"`
}

// ToMyMailboxRequest converts the InboxSearch to a MyMailboxRequest
func (i *InboxSearch) ToMyMailboxRequest() *api.MyMailboxRequest {

	page := i.Page
	if i.Page == nil {
		*page = 1
	}
	filters := make(map[string]interface{})
	if i.Filter != nil {
		if i.Filter.Subject != nil {
			filters["subject"] = *i.Filter.Subject
		}
		if i.Filter.State != nil {
			filters["state"] = i.Filter.State
		}
		if i.Filter.Type != nil {
			filters["type"] = *i.Filter.Type
		}
	}

	return &api.MyMailboxRequest{
		Filter: &filters,
		Page:   page,
	}
}

type InboxSearchResponse struct {
	// Total The total number of results
	Total *float32 `json:"total,omitempty"`
	// Results The results of the search
	Results []*InboxSearchResult `json:"results,omitempty"`
}

type InboxSearchResult struct {
	*api.Item
}

type OutboxFilter struct {
	Subject *string   `json:"subject,omitempty"`
	State   []*string `json:"state,omitempty"`
	Type    *string   `json:"type,omitempty"`
}

// ToOutboxRequest converts the OutboxSearch to an OutboxRequest
func (o *OutboxSearch) ToOutboxRequest(tenant *string) *api.MyOutboxRequest {
	page := o.Page
	if o.Page == nil {
		*page = 1
	}
	filters := make(map[string]interface{})
	if o.Filter != nil {
		if o.Filter.Subject != nil {
			filters["subject"] = *o.Filter.Subject
		}
		if o.Filter.State != nil {
			filters["state"] = o.Filter.State
		}
		if o.Filter.Type != nil {
			filters["type"] = *o.Filter.Type
		}
	}

	return &api.MyOutboxRequest{
		Filter:     &filters,
		Page:       page,
		SenderUser: o.SenderUser,
	}
}

type OutboxSearch struct {
	// Filter The filter to apply to the search
	Filter *OutboxFilter `json:"filter,omitempty"`
	// Page The page number to get
	Page *float32 `json:"page,omitempty"`
	// The user id or api key of the sender. By default the logged in user is used.
	SenderUser *string `json:"sender_user,omitempty"`
}

type OutboxSearchResponse struct {
	// Total The total number of results
	Total *float32 `json:"total,omitempty"`
	// Results The results of the search
	Results []*OutboxSearchResult `json:"results,omitempty"`
}
type OutboxSearchResult struct {
	*api.Item
}
