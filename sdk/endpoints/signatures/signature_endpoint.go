package signatures

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// ExportSignature exports a signature by its ID in the specified format.
func ExportSignature(client *sdkClient.BrifleClient, ctx context.Context, signatureId *string, exportOptions *ExportOptions) (*string, *api.ResponseStatus, error) {
	if signatureId == nil {
		return nil, nil, errors.New("signatureId is required")
	}

	if exportOptions == nil {
		return nil, nil, errors.New("exportOptions are required")
	}

	var format api.WebApiControllerSignatureControllerExportSignatureParamsFormat
	switch strings.ToLower(exportOptions.Format) {
	case "xml":
		format = api.Xml
	default:
		format = api.Xml
	}
	exportResponse, err := client.ApiClient.WebApiControllerSignatureControllerExportSignature(ctx, *signatureId, format)
	if err != nil {
		return nil, nil, err
	}
	if exportResponse == nil {
		return nil, nil, errors.New("export response is nil")
	}
	status, res, err := api.ParseResponseAsString(exportResponse)
	if err != nil {
		return nil, nil, err
	}
	return &res, status, nil
}

// CreateSignatureReference creates a new signature reference with the provided options.
func CreateSignatureReference(client *sdkClient.BrifleClient, ctx context.Context, tenantId *string, signatureReferenceOptions *SignatureReferenceOptions) (*SignatureReference, *api.ResponseStatus, error) {
	if tenantId == nil {
		return nil, nil, errors.New("tenantId is required")
	}

	if signatureReferenceOptions == nil {
		return nil, nil, errors.New("signatureReferenceOptions are required")
	}

	request := signatureReferenceOptions.ToApiSignatureReferenceOptions()
	if request == nil {
		return nil, nil, errors.New("signatureReferenceOptions are invalid")
	}
	resp, err := client.ApiClient.WebApiControllerSignatureControllerCreateSignatureReference(ctx, *tenantId, *request)
	if err != nil {
		return nil, nil, err
	}
	if resp == nil {
		return nil, nil, errors.New("response is nil")
	}

	var response *api.SignatureReference
	status, err := api.ValidateHttpResponse(err, resp, &response)
	if err != nil {
		return nil, status, err
	}
	if status != nil && status.HttpStatus != 200 {
		return nil, status, fmt.Errorf("unexpected HTTP status: %d", status.HttpStatus)
	}

	signatureReference := &SignatureReference{}
	err = signatureReference.ParseFromApiResponse(response)
	if err != nil {
		return nil, nil, err
	}

	return signatureReference, status, nil
}

type SignatureReferenceOptions struct {
	Fields []SignatureReferenceField `json:"fields"`
}

func (s *SignatureReferenceOptions) ToApiSignatureReferenceOptions() *api.CreateSignatureReferenceRequest {
	if s == nil {
		return nil
	}

	var fields []struct {
		// Name The name of the field
		Name string `json:"name"`

		// Purpose The purpose of the field
		Purpose *string `json:"purpose,omitempty"`

		// Role The role of the signer of the field
		Role *string `json:"role,omitempty"`
	}
	if s.Fields != nil {
		fields = make([]struct {
			// Name The name of the field
			Name string `json:"name"`
			// Purpose The purpose of the field
			Purpose *string `json:"purpose,omitempty"`
			// Role The role of the signer of the field
			Role *string `json:"role,omitempty"`
		}, len(s.Fields))
		for i, field := range s.Fields {
			fields[i] = struct {
				// Name The name of the field
				Name string `json:"name"`
				// Purpose The purpose of the field
				Purpose *string `json:"purpose,omitempty"`
				// Role The role of the signer of the field
				Role *string `json:"role,omitempty"`
			}{
				Name:    field.Name,
				Purpose: &field.Purpose,
				Role:    &field.Role,
			}
		}
	}
	return &api.CreateSignatureReferenceRequest{
		Fields: &fields,
	}
}

type SignatureReferenceField struct {
	Name    string `json:"name"`
	Purpose string `json:"purpose"`
	Role    string `json:"role"`
}

type SignatureReference struct {
	DocumentSignatures []string                  `json:"document_signatures,omitempty"` // JSON array of strings
	Id                 string                    `json:"id"`
	ManagedBy          string                    `json:"managed_by"`                  // The ID of the user who manages the signature reference
	SignaturesFields   []SignatureReferenceField `json:"signatures_fields,omitempty"` // JSON array of SignatureReferenceField
}

// ParseFromApiResponse parses the SignatureReference from the API response.
func (signatureReference *SignatureReference) ParseFromApiResponse(response *api.SignatureReference) error {
	if response == nil {
		return errors.New("response is nil")
	}

	// Parse JSON array to a slice of strings
	if response.DocumentSignatures != nil {
		signatureReference.DocumentSignatures = make([]string, 0)
		var raw string
		err := json.Unmarshal([]byte(response.DocumentSignatures), &raw)
		if err != nil {
			return errors.New("failed to unmarshal document signatures: " + err.Error())
		}
		err = json.Unmarshal([]byte(raw), &signatureReference.DocumentSignatures)
		if err != nil {
			return errors.New("failed to unmarshal document signatures: " + err.Error())
		}
	}

	signatureReference.Id = *response.Id
	signatureReference.ManagedBy = *response.ManagedBy

	if response.SignatureFields != nil {
		var raw string
		err := json.Unmarshal([]byte(response.SignatureFields), &raw)
		if err != nil {
			return errors.New("failed to unmarshal signature fields: " + err.Error())
		}
		signatureReference.SignaturesFields = make([]SignatureReferenceField, 0)
		err = json.Unmarshal([]byte(raw), &signatureReference.SignaturesFields)
		if err != nil {
			return errors.New("failed to unmarshal signature fields: " + err.Error())
		}
	}

	return nil
}

type ExportOptions struct {
	Format string `json:"format"`
}
