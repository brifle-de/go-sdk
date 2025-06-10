package content

import (
	"context"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

func buildReceiver(receiver *ReceiverData) *api.ApiSendContentReceiverRequest {
	if receiver == nil {
		return nil
	}

	if receiver.BirthInformation != nil {
		return &api.ApiSendContentReceiverRequest{
			BirthInformation: &api.ApiSendContentReceiverBirthInformation{
				GivenNames:    receiver.BirthInformation.FirstName,
				LastName:      receiver.BirthInformation.LastName,
				PlaceOfBirth:  receiver.BirthInformation.PlaceOfBirth,
				DateOfBirth:   receiver.BirthInformation.DateOfBirth,
				PostalAddress: receiver.BirthInformation.PostalAddress,
				BirthName:     receiver.BirthInformation.NameAtBirth,
			},
		}
	}

	if receiver.Email != nil {
		return &api.ApiSendContentReceiverRequest{
			Email:       receiver.Email.Email,
			FullName:    receiver.Email.Name,
			DateOfBirth: receiver.Email.DateOfBirth,
		}
	}

	if receiver.Phone != nil {
		return &api.ApiSendContentReceiverRequest{
			Tel:         receiver.Phone.PhoneNumber,
			FullName:    receiver.Phone.Name,
			DateOfBirth: receiver.Phone.DateOfBirth,
		}
	}

	return nil
}

// SendContent sends a document to the specified receiver with the provided content and metadata.
func SendContent(client *sdkClient.BrifleClient, context context.Context, tenant *string, sendContent *SendContentRequest) (*SendDocumentResponse, *api.ResponseStatus, error) {
	if sendContent == nil {
		return nil, nil, errors.New("send content request is nil")
	}

	receiver := buildReceiver(sendContent.To)
	if receiver == nil {
		return nil, nil, errors.New("receiver data is invalid")
	}

	convertedBody := make([]api.ApiSendContentContentRequest, len(*sendContent.Body))
	for i, item := range *sendContent.Body {
		convertedBody[i] = api.ApiSendContentContentRequest{
			Content: item.Content,
			Type:    item.Type,
		}
	}

	request := &api.ApiSendContentSendContentRequest{
		To:            receiver,
		Type:          (*api.ApiSendContentSendContentRequestType)(sendContent.Type),
		Body:          &convertedBody,
		Subject:       sendContent.Subject,
		PaymentInfo:   sendContent.PaymentInfo.ToApiPaymentInfo(),
		SignatureInfo: sendContent.SignatureInfo.ToApiSignatureInfo(),
	}

	response, err := client.ApiClient.WebApiControllerContentControllerSend(context, *tenant, *request)
	if err != nil {
		return nil, nil, err
	}

	var res SendDocumentResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}

	return &res, status, nil
}

// CheckReceiver checks if the receiver data is valid and returns a response indicating the result.
func CheckReceiver(client *sdkClient.BrifleClient, context context.Context, receiver *ReceiverData) (*ReceiverCheckResponse, *api.ResponseStatus, error) {
	if receiver == nil {
		return nil, nil, errors.New("receiver data is nil")
	}
	r := buildReceiver(receiver)
	response, err := client.ApiClient.WebApiControllerContentControllerCheckReceiver(context, *r)
	var res ReceiverCheckResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// GetContent retrieves the content of a document by its ID. If readFlag is set to true, it marks the document as read.
func GetContent(client *sdkClient.BrifleClient, context context.Context, documentId *string, readFlag *bool) (*DocumentResponse, *api.ResponseStatus, error) {
	params := api.WebApiControllerContentControllerGetParams{
		Read: readFlag, // Set to true if you want to mark the document as read
	}
	response, err := client.ApiClient.WebApiControllerContentControllerGet(context, *documentId, &params)
	var res DocumentResponse
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// Types

type ReceiverData struct {
	BirthInformation *BirthInformationReceiver `json:"birth_information,omitempty"`
	Email            *EmailReceiver            `json:"email,omitempty"`
	Phone            *PhoneReceiver            `json:"phone,omitempty"`
}

type BirthInformationReceiver struct {
	FirstName     *string `json:"first_name,omitempty"`
	LastName      *string `json:"last_name,omitempty"`
	PlaceOfBirth  *string `json:"place_of_birth,omitempty"`
	DateOfBirth   *string `json:"date_of_birth,omitempty"`
	PostalAddress *string `json:"postal_address,omitempty"`
	NameAtBirth   *string `json:"name_at_birth,omitempty"`
}
type EmailReceiver struct {
	Email       *string `json:"email,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Name        *string `json:"name,omitempty"`
}
type PhoneReceiver struct {
	PhoneNumber *string `json:"phone_number,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Name        *string `json:"name,omitempty"`
}

type DocumentResponse struct {
	*api.ContentGetResponse
}

type ReceiverCheckResponse struct {
	*api.ReceiverExistResponse
}

type SendDocumentResponse struct {
	*api.ContentCreateResponse
}

type SendContentRequest struct {
	To            *ReceiverData  `json:"to,omitempty"`
	Type          *string        `json:"type,omitempty"`
	Body          *[]ContentItem `json:"body,omitempty"`
	Subject       *string        `json:"subject,omitempty"`
	PaymentInfo   *PaymentInfo   `json:"payment_info,omitempty"`
	SignatureInfo *SignatureInfo `json:"signature_info,omitempty"`
}

type ContentItem struct {
	// Content Content
	Content *string `json:"content,omitempty"`

	// Type Type
	Type *string `json:"type,omitempty"`
}

type PaymentDetails struct {
	// Amount Amount
	Amount *float32 `json:"amount,omitempty"`

	// Currency Currency
	Currency *string `json:"currency,omitempty"`

	// Description Description
	Description *string `json:"description,omitempty"`

	// DueDate Due Date
	DueDate *string `json:"due_date,omitempty"`

	// Iban IBAN
	Iban *string `json:"iban,omitempty"`

	// Reference Reference
	Reference *string `json:"reference,omitempty"`
}

type PaymentInfo struct {
	Details *PaymentDetails `json:"details,omitempty"`

	// Payable Payable
	Payable *bool `json:"payable,omitempty"`
}

type Signer struct {
}

type SignatureInfo struct {
	RequestingSigner *[]struct {
		// Field Field
		Field *string `json:"field,omitempty"`

		// Signer Which party shall sign the field
		Signer *string `json:"signer,omitempty"`
	} `json:"requesting_signer,omitempty"`

	// SignatureReference Signature Reference
	SignatureReference *string `json:"signature_reference,omitempty"`
}

func (reqSigner *SignatureInfo) ToApiSignatureInfo() *api.ApiSendContentSignatureInfo {
	if reqSigner == nil {
		return nil
	}
	if reqSigner.RequestingSigner == nil {
		return nil
	}
	convertedRequestingSigner := make([]struct {
		Field  *string                                                `json:"field,omitempty"`
		Signer *api.ApiSendContentSignatureInfoRequestingSignerSigner `json:"signer,omitempty"`
	}, len(*reqSigner.RequestingSigner))
	for i, item := range *reqSigner.RequestingSigner {
		signer := api.ApiSendContentSignatureInfoRequestingSignerSigner(*item.Signer)
		convertedRequestingSigner[i] = struct {
			Field  *string                                                `json:"field,omitempty"`
			Signer *api.ApiSendContentSignatureInfoRequestingSignerSigner `json:"signer,omitempty"`
		}{
			Field:  item.Field,
			Signer: &signer,
		}
	}
	return &api.ApiSendContentSignatureInfo{
		RequestingSigner:   &convertedRequestingSigner,
		SignatureReference: reqSigner.SignatureReference,
	}
}

func (paymentInfo *PaymentInfo) ToApiPaymentInfo() *api.ApiSendContentPaymentInfo {
	if paymentInfo == nil {
		return nil
	}
	if paymentInfo.Details == nil {
		return &api.ApiSendContentPaymentInfo{
			Payable: paymentInfo.Payable,
		}
	}
	return &api.ApiSendContentPaymentInfo{
		Details: &api.ApiSendContentPaymentDetails{
			Amount:      paymentInfo.Details.Amount,
			Currency:    paymentInfo.Details.Currency,
			Description: paymentInfo.Details.Description,
			DueDate:     paymentInfo.Details.DueDate,
			Iban:        paymentInfo.Details.Iban,
			Reference:   paymentInfo.Details.Reference,
		},
		Payable: paymentInfo.Payable,
	}
}

const (
	Letter   = "letter"
	Invoice  = "invoice"
	Contract = "contract"
)
