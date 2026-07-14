package content

import (
	"context"
	"errors"
	"strings"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// Cover letter output formats used by GetCoverLetter.
const (
	FormatPdf    = "pdf"
	FormatBase64 = "base64"
)

// UploadCoverLetter uploads a custom cover letter template for the given tenant.
// The content must be base64 encoded (see sdk.Base64Encode).
func UploadCoverLetter(client *sdkClient.BrifleClient, ctx context.Context, tenant *string, name *string, contentBase64 *string) (*CoverLetter, *api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, nil, errors.New("tenant is required")
	}
	if name == nil || *name == "" {
		return nil, nil, errors.New("name is required")
	}
	if contentBase64 == nil || *contentBase64 == "" {
		return nil, nil, errors.New("content is required")
	}
	request := api.UpdateCoverLetterRequest{
		Name:    *name,
		Content: *contentBase64,
	}
	response, err := client.ApiClient.WebApiControllerContentControllerUploadCoverLetter(ctx, *tenant, request)
	var res CoverLetter
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// DeleteCoverLetter deletes a custom cover letter template by name.
func DeleteCoverLetter(client *sdkClient.BrifleClient, ctx context.Context, tenant *string, name *string) (*api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, errors.New("tenant is required")
	}
	if name == nil || *name == "" {
		return nil, errors.New("name is required")
	}
	response, err := client.ApiClient.WebApiControllerContentControllerDeleteCoverLetter(ctx, *tenant, *name)
	if err != nil {
		return nil, err
	}
	status, _, err := api.ParseResponseAsString(response)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// ListCoverLetters lists the cover letters available to the tenant, both the
// built-in "default" templates and the tenant's "custom" templates.
func ListCoverLetters(client *sdkClient.BrifleClient, ctx context.Context, tenant *string) (*CoverLettersList, *api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, nil, errors.New("tenant is required")
	}
	response, err := client.ApiClient.WebApiControllerContentControllerGetCoverLettersList(ctx, *tenant)
	var res CoverLettersList
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// GetCoverLetter retrieves the content of a cover letter. coverType is
// CoverLetterDefault or CoverLetterCustom; format is FormatPdf or FormatBase64.
// The content is returned as raw bytes (a PDF, or base64 text when using
// FormatBase64).
func GetCoverLetter(client *sdkClient.BrifleClient, ctx context.Context, tenant *string, coverType *string, fileName *string, format *string) ([]byte, *api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, nil, errors.New("tenant is required")
	}
	if coverType == nil || *coverType == "" {
		return nil, nil, errors.New("cover letter type is required")
	}
	if fileName == nil || *fileName == "" {
		return nil, nil, errors.New("file name is required")
	}

	var apiFormat api.WebApiControllerContentControllerGetCoverLetterParamsFormat
	switch strings.ToLower(strVal(format)) {
	case FormatBase64:
		apiFormat = api.Base64
	default:
		apiFormat = api.Pdf
	}

	response, err := client.ApiClient.WebApiControllerContentControllerGetCoverLetter(ctx, *tenant, *coverType, *fileName, apiFormat)
	if err != nil {
		return nil, nil, err
	}
	status, body, err := api.ParseResponseAsBytes(response)
	if err != nil {
		return nil, status, err
	}
	return body, status, nil
}

// CoverLetter is the metadata of an uploaded cover letter.
type CoverLetter struct {
	*api.UpdateCoverLetterResponse
}

// CoverLettersList is the list of cover letters available to a tenant.
type CoverLettersList struct {
	*api.CoverLettersOverviewResponse
}
