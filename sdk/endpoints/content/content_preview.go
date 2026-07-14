package content

import (
	"context"
	"errors"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// Cover letter types used by paper mail previews and cover letter management.
const (
	CoverLetterDefault = "default"
	CoverLetterCustom  = "custom"
)

// PreviewPaperMail renders how a document would look when delivered as physical
// paper mail and returns the resulting PDF as raw bytes.
func PreviewPaperMail(client *sdkClient.BrifleClient, ctx context.Context, tenant *string, req *PreviewPaperMailRequest) ([]byte, *api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, nil, errors.New("tenant is required")
	}
	if req == nil {
		return nil, nil, errors.New("preview paper mail request is nil")
	}
	if req.To == nil || req.Body == nil || req.CoverLetter == nil {
		return nil, nil, errors.New("to, body and cover_letter are required")
	}

	apiReq := api.SendContentPreviewPapermailRequest{
		Body: api.Body{
			Content: strVal(req.Body.Content),
			Type:    strVal(req.Body.Type),
		},
		CoverLetter: api.PreviewPapermailCoverLetter{
			Enable: req.CoverLetter.Enable,
			Name:   strVal(req.CoverLetter.Name),
			Data:   req.CoverLetter.Data,
			Type:   api.PreviewPapermailCoverLetterType(strVal(req.CoverLetter.Type)),
		},
		To: api.PreviewPapermailReceiver{
			AddressLine1: strVal(req.To.AddressLine1),
			AddressLine2: req.To.AddressLine2,
			City:         strVal(req.To.City),
			Country:      strVal(req.To.Country),
			PostalCode:   strVal(req.To.PostalCode),
		},
	}

	response, err := client.ApiClient.WebApiControllerContentControllerPreviewPaperMail(ctx, *tenant, apiReq)
	if err != nil {
		return nil, nil, err
	}
	status, body, err := api.ParseResponseAsBytes(response)
	if err != nil {
		return nil, status, err
	}
	return body, status, nil
}

// PreviewPaperMailRequest describes a paper mail preview.
type PreviewPaperMailRequest struct {
	To          *PreviewReceiver    `json:"to,omitempty"`
	CoverLetter *PreviewCoverLetter `json:"cover_letter,omitempty"`
	Body        *PreviewBody        `json:"body,omitempty"`
}

// PreviewReceiver is the recipient address printed on the paper mail.
type PreviewReceiver struct {
	AddressLine1 *string `json:"address_line1,omitempty"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         *string `json:"city,omitempty"`
	Country      *string `json:"country,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
}

// PreviewCoverLetter configures the cover letter of the paper mail.
type PreviewCoverLetter struct {
	// Enable controls whether a cover letter is included.
	Enable bool `json:"enable"`
	// Name of the cover letter template. Ignored when Data is provided.
	Name *string `json:"name,omitempty"`
	// Type of the cover letter: CoverLetterDefault or CoverLetterCustom.
	Type *string `json:"type,omitempty"`
	// Data is custom cover letter content. Optional.
	Data *string `json:"data,omitempty"`
}

// PreviewBody is the document content to render as paper mail.
type PreviewBody struct {
	// Content is the base64 encoded document.
	Content *string `json:"content,omitempty"`
	// Type is the MIME type, e.g. "application/pdf".
	Type *string `json:"type,omitempty"`
}
