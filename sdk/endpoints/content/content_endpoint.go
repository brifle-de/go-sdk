package content

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// strVal safely dereferences a *string, returning "" when nil.
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// f32Val safely dereferences a *float32, returning 0 when nil.
func f32Val(f *float32) float32 {
	if f == nil {
		return 0
	}
	return *f
}

func buildReceiver(receiver *ReceiverData) *api.ApiSendContentReceiverRequest {
	if receiver == nil {
		return nil
	}

	if receiver.BirthInformation != nil {
		return &api.ApiSendContentReceiverRequest{
			BirthInformation: &api.ApiSendContentReceiverBirthInformation{
				GivenNames:    strVal(receiver.BirthInformation.FirstName),
				LastName:      strVal(receiver.BirthInformation.LastName),
				PlaceOfBirth:  strVal(receiver.BirthInformation.PlaceOfBirth),
				DateOfBirth:   strVal(receiver.BirthInformation.DateOfBirth),
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

// SendContent sends a document to the specified receiver with the provided
// content and metadata, on behalf of the given tenant.
//
//	req := content.SendContentRequest{
//		To: &content.ReceiverData{
//			Email: &content.EmailReceiver{
//				Email: sdk.String("max@example.com"),
//				Name:  sdk.String("Max Mustermann"),
//			},
//		},
//		Type:    sdk.String(content.Letter),
//		Subject: sdk.String("Welcome to Brifle"),
//		Body: &[]content.ContentItem{
//			{Content: sdk.Base64Encode(pdfBytes), Type: sdk.String("application/pdf")},
//		},
//	}
//	res, respStatus, err := content.SendContent(client, ctx, &tenant, &req)
//	if err == nil && respStatus.HttpStatus == 200 {
//		fmt.Println("document id:", *res.Id)
//	}
//
// See [ReceiverData] for the ways to address a recipient, and the [Letter],
// [Invoice] and [Contract] constants for the document type.
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
			Type:    (*api.ApiSendContentContentRequestType)(item.Type),
		}
	}

	var contentType api.ApiSendContentSendContentRequestType
	if sendContent.Type != nil {
		contentType = api.ApiSendContentSendContentRequestType(*sendContent.Type)
	}

	request := &api.ApiSendContentSendContentRequest{
		To:            *receiver,
		Type:          contentType,
		Body:          convertedBody,
		Subject:       strVal(sendContent.Subject),
		Fallback:      sendContent.Fallback.ToApiFallback(),
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

// GetContentAction retrieves the actions available for a document by its ID. If readFlag is set to true, it marks the document as read.
func GetContentAction(client *sdkClient.BrifleClient, context context.Context, documentId *string) (*ContentActions, *api.ResponseStatus, error) {
	response, err := client.ApiClient.WebApiControllerContentControllerGetActions(context, *documentId)
	var res ContentActions
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// GetDeliveryCertificate retrieves the delivery certificate for a document by its ID.
func GetDeliveryCertificate(client *sdkClient.BrifleClient, context context.Context, documentId *string) (*DeliveryCertificate, *api.ResponseStatus, error) {
	if documentId == nil || *documentId == "" {
		return nil, nil, errors.New("document ID is required")
	}

	response, err := client.ApiClient.WebApiControllerContentControllerGetDeliveryCertificate(context, *documentId)
	var res DeliveryCertificate
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// GetDeliveryCertificate retrieves the delivery certificate for a document by its ID.
type DeliveryCertificate struct {
	// Certificate The delivery certificate in XML format.
	Certificate *string `json:"certificate,omitempty"`
	Meta        *struct {
		// DocumentId The ID of the document for which the delivery certificate is requested.
		DocumentId *string `json:"document_id,omitempty"`

		// Id The ID of the delivery certificate.
		Id *string `json:"id,omitempty"`

		// Type The type of the delivery certificate. aes - advanced electronic seal.
		Type *string `json:"type,omitempty"`
	} `json:"meta,omitempty"`
}

type ContentActions struct {
	Payments *struct {
		Details *struct {
			// Amount the amount to pay in the smallest unit of the currency
			Amount *float32 `json:"amount,omitempty"`

			// Currency the currency of the payment
			Currency *string `json:"currency,omitempty"`

			// Iban the iban of the payment
			Iban *string `json:"iban,omitempty"`

			// Market the market of the payment. Important for the payment provider, e.g. Tink
			Market *string `json:"market,omitempty"`

			// Reference the reference of the payment
			Reference *string `json:"reference,omitempty"`

			// TinkPaymentId the payment id in the Tink system
			TinkPaymentId *string `json:"tink_payment_id,omitempty"`
		} `json:"details,omitempty"`
		Link *string `json:"link,omitempty"`
	} `json:"payments,omitempty"`
	Signatures *struct {
		DocumentSignatures *struct {
			// SignatureIds array of signature ids
			SignatureIds       json.RawMessage `json:"signature_ids,omitempty"`
			SignatureReference *string         `json:"signature_reference,omitempty"`
		} `json:"document_signatures,omitempty"`
		EmbeddedSignatures *[]struct {
			Id                  *string `json:"id,omitempty"`
			CreatedBy           *string `json:"created_by,omitempty"`
			CreatedDate         *string `json:"created_date,omitempty"`
			DocumentSignatureId *string `json:"document_signature_id,omitempty"`
			DueDate             *string `json:"due_date,omitempty"`
			FieldName           *string `json:"field_name,omitempty"`
			History             *string `json:"history,omitempty"`

			// Purpose the purpose why the signature is needed. Important if a document requires multiple signatures from the same signer
			Purpose       *string `json:"purpose,omitempty"`
			RequestDate   *string `json:"request_date,omitempty"`
			RequestedTo   *string `json:"requested_to,omitempty"`
			SignatureDate *string `json:"signature_date,omitempty"`
			SignedBy      *string `json:"signed_by,omitempty"`
			SignedFor     *string `json:"signed_for,omitempty"`
			Value         *string `json:"value,omitempty"`
		} `json:"embedded_signatures,omitempty"`
		SignatureReference *struct {
			DocumentSignatures *string `json:"document_signatures,omitempty"`
			ManagedBy          *string `json:"managed_by,omitempty"`
			SignatureFields    *string `json:"signature_fields,omitempty"`
		} `json:"signature_reference,omitempty"`
	} `json:"signatures,omitempty"`
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
	Fallback      *Fallback      `json:"fallback,omitempty"`
}

type Fallback struct {
	// EnabledPhysicalDelivery Enable Physical Delivery
	EnabledPhysicalDelivery bool       `json:"enabled_physical_delivery,omitempty"`
	PaperMail               *PaperMail `json:"paper_mail,omitempty"`
}

// ToApiFallback converts the Fallback to an ApiSendContentFallback
func (f *Fallback) ToApiFallback() *api.ApiSendContentFallback {
	if f == nil {
		return nil
	}

	enabled := f.EnabledPhysicalDelivery

	if f.PaperMail == nil || f.PaperMail.Recipient == nil {
		return &api.ApiSendContentFallback{
			EnabledPhysicalDelivery: &enabled,
		}
	}

	r := f.PaperMail.Recipient
	return &api.ApiSendContentFallback{
		EnabledPhysicalDelivery: &enabled,
		PaperMail: &struct {
			Recipient *struct {
				// AddressLine1 Address Line 1
				AddressLine1 string `json:"address_line1"`

				// AddressLine2 Address Line 2
				AddressLine2 *string `json:"address_line2,omitempty"`

				// AddressLine3 Address Line 3
				AddressLine3 *string `json:"address_line3,omitempty"`

				// City City
				City string `json:"city"`

				// Country Country, the country code in ISO 3166-1 alpha-2 format
				Country *string `json:"country,omitempty"`

				// PostalCode ZIP Code
				PostalCode string `json:"postal_code"`
			} `json:"recipient,omitempty"`
		}{
			Recipient: &struct {
				// AddressLine1 Address Line 1
				AddressLine1 string `json:"address_line1"`

				// AddressLine2 Address Line 2
				AddressLine2 *string `json:"address_line2,omitempty"`

				// AddressLine3 Address Line 3
				AddressLine3 *string `json:"address_line3,omitempty"`

				// City City
				City string `json:"city"`

				// Country Country, the country code in ISO 3166-1 alpha-2 format
				Country *string `json:"country,omitempty"`

				// PostalCode ZIP Code
				PostalCode string `json:"postal_code"`
			}{
				AddressLine1: strVal(r.AddressLine1),
				AddressLine2: r.AddressLine2,
				AddressLine3: r.AddressLine3,
				City:         strVal(r.City),
				Country:      r.Country,
				PostalCode:   strVal(r.PostalCode),
			},
		},
	}
}

type PaperMail struct {
	Recipient *Recipient `json:"recipient,omitempty"`
}

type Recipient struct {
	AddressLine1 *string `json:"address_line1,omitempty"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	AddressLine3 *string `json:"address_line3,omitempty"`
	City         *string `json:"city,omitempty"`
	Country      *string `json:"country,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
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
		// Field Field
		Field string `json:"field"`

		// Signer Which party shall sign the field
		Signer api.ApiSendContentSignatureInfoRequestingSignerSigner `json:"signer"`
	}, len(*reqSigner.RequestingSigner))
	for i, item := range *reqSigner.RequestingSigner {
		var signer api.ApiSendContentSignatureInfoRequestingSignerSigner
		if item.Signer != nil {
			signer = api.ApiSendContentSignatureInfoRequestingSignerSigner(*item.Signer)
		}
		convertedRequestingSigner[i] = struct {
			// Field Field
			Field string `json:"field"`

			// Signer Which party shall sign the field
			Signer api.ApiSendContentSignatureInfoRequestingSignerSigner `json:"signer"`
		}{
			Field:  strVal(item.Field),
			Signer: signer,
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
			Amount:      f32Val(paymentInfo.Details.Amount),
			Currency:    strVal(paymentInfo.Details.Currency),
			Description: strVal(paymentInfo.Details.Description),
			DueDate:     strVal(paymentInfo.Details.DueDate),
			Iban:        strVal(paymentInfo.Details.Iban),
			Reference:   strVal(paymentInfo.Details.Reference),
		},
		Payable: paymentInfo.Payable,
	}
}

const (
	Letter   = "letter"
	Invoice  = "invoice"
	Contract = "contract"
)
