package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/brifle-de/brifle-sdk/sdk/api"
	sdkClient "github.com/brifle-de/brifle-sdk/sdk/client"
)

// Wallet item types.
const (
	ProofOfOwnership  = "proof_of_ownership"
	ProofOfPermission = "proof_of_permission"
)

// Data element types.
const (
	TypeText = "text"
	TypeUrl  = "url"
	TypeJwt  = "jwt"
)

// CreateWalletItem issues a new wallet item for the given tenant. On success the
// API responds with HTTP 201 and the metadata of the created item.
//
//	req := wallet.CreateWalletRequest{
//		Type:    wallet.ProofOfOwnership,
//		Subject: "Membership Card",
//		Data: []wallet.DataElement{{
//			Name:        sdk.String("Member Name"),
//			Value:       sdk.String("Max Mustermann"),
//			Type:        sdk.String(wallet.TypeText),
//			ReferenceId: sdk.String("member_name"),
//		}},
//	}
//	res, respStatus, err := wallet.CreateWalletItem(client, ctx, &tenant, &req)
//	if err == nil && respStatus.HttpStatus == 201 {
//		fmt.Println("created wallet item:", *res.Id)
//	}
func CreateWalletItem(client *sdkClient.BrifleClient, ctx context.Context, tenant *string, req *CreateWalletRequest) (*WalletMeta, *api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, nil, errors.New("tenant is required")
	}
	if req == nil {
		return nil, nil, errors.New("create wallet request is nil")
	}
	apiReq := req.toApiRequest()
	response, err := client.ApiClient.WebApiControllerWalletControllerCreateWalletItem(ctx, *tenant, *apiReq)
	var res WalletMeta
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// ReadWalletItem retrieves a wallet item (its data and metadata) by its ID.
func ReadWalletItem(client *sdkClient.BrifleClient, ctx context.Context, tenant *string, id *string) (*WalletItem, *api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, nil, errors.New("tenant is required")
	}
	if id == nil || *id == "" {
		return nil, nil, errors.New("id is required")
	}
	response, err := client.ApiClient.WebApiControllerWalletControllerReadWalletItem(ctx, *tenant, *id)
	var res WalletItem
	status, err := api.ValidateHttpResponse(err, response, &res)
	if err != nil {
		return nil, status, err
	}
	return &res, status, nil
}

// RevokeWalletItem revokes (deletes) a wallet item by its ID. An item can only
// be revoked while it has not yet been assigned to a user wallet.
func RevokeWalletItem(client *sdkClient.BrifleClient, ctx context.Context, tenant *string, id *string) (*api.ResponseStatus, error) {
	if tenant == nil || *tenant == "" {
		return nil, errors.New("tenant is required")
	}
	if id == nil || *id == "" {
		return nil, errors.New("id is required")
	}
	response, err := client.ApiClient.WebApiControllerWalletControllerDeleteWalletItem(ctx, *tenant, *id)
	if err != nil {
		return nil, err
	}
	status, _, err := api.ParseResponseAsString(response)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// WalletMeta is the metadata of a wallet item.
type WalletMeta struct {
	*api.WalletMetaResponse
}

// WalletItem is a full wallet item (data + metadata).
type WalletItem struct {
	*api.WalletItemResponse
}

// CreateWalletRequest is the friendly request used to create a wallet item.
type CreateWalletRequest struct {
	// Type of the wallet item: ProofOfOwnership or ProofOfPermission.
	Type string `json:"type"`
	// Subject of the wallet item, e.g. "Driver's License".
	Subject string `json:"subject"`
	// Data elements contained in the wallet item.
	Data []DataElement `json:"data"`
	// Styles for how the wallet item is displayed. Optional.
	Styles *WalletStyle `json:"styles,omitempty"`
	// ExportWallet controls export to external wallets (Apple/Google). Optional.
	ExportWallet *ExportWallet `json:"export_wallet,omitempty"`
	// Immutable, if true, means the item cannot be changed once created.
	Immutable *bool `json:"immutable,omitempty"`
	// RetentionPeriodDays to retain an unassigned item (defaults to 30).
	RetentionPeriodDays *int `json:"retention_period_days,omitempty"`
	// ExpiresAt is the expiration time; nil means no expiration.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// NotBefore is the time before which the item is not valid. Optional.
	NotBefore *time.Time `json:"not_before,omitempty"`
}

// DataElement is a single piece of data stored in a wallet item.
type DataElement struct {
	// Name of the data element.
	Name *string `json:"name,omitempty"`
	// Value of the data element.
	Value *string `json:"value,omitempty"`
	// Type of the data element: TypeText, TypeUrl or TypeJwt.
	Type *string `json:"type,omitempty"`
	// ReferenceId uniquely identifies the element within the item (required for
	// optimised export). It must not contain sensitive information.
	ReferenceId *string `json:"reference_id,omitempty"`
}

// WalletStyle configures the visual presentation of a wallet item.
type WalletStyle struct {
	BackgroundColor   *string `json:"background_color,omitempty"`
	OnBackgroundColor *string `json:"on_background_color,omitempty"`
	PrimaryColor      *string `json:"primary_color,omitempty"`
	OnPrimaryColor    *string `json:"on_primary_color,omitempty"`
	SecondaryColor    *string `json:"secondary_color,omitempty"`
	OnSecondaryColor  *string `json:"on_secondary_color,omitempty"`
	PrimaryQrCodeData *string `json:"primary_qr_code_data,omitempty"`
	Rows              []Row   `json:"rows,omitempty"`
}

// Row displays two referenced data elements side by side in a wallet item.
type Row struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

// ExportWallet controls exporting a wallet item to external wallets.
type ExportWallet struct {
	AppleWallet  *bool `json:"apple_wallet,omitempty"`
	GoogleWallet *bool `json:"google_wallet,omitempty"`
}

func (r *CreateWalletRequest) toApiRequest() *api.CreateWalletRequest {
	if r == nil {
		return nil
	}

	elements := make([]api.DataElement, len(r.Data))
	for i, d := range r.Data {
		elements[i] = api.DataElement{
			Name:        d.Name,
			Value:       d.Value,
			ReferenceId: d.ReferenceId,
			Type:        (*api.DataElementType)(d.Type),
		}
	}

	req := &api.CreateWalletRequest{
		Type:                api.CreateWalletRequestType(r.Type),
		Subject:             r.Subject,
		Data:                api.WalletDataElement{Elements: elements},
		Immutable:           r.Immutable,
		RetentionPeriodDays: r.RetentionPeriodDays,
		ExpiresAt:           r.ExpiresAt,
		NotBefore:           r.NotBefore,
	}

	if r.ExportWallet != nil {
		req.ExportWallet = &api.ExportWallet{
			AppleWallet:  r.ExportWallet.AppleWallet,
			GoogleWallet: r.ExportWallet.GoogleWallet,
		}
	}

	if r.Styles != nil {
		style := &api.WalletStyle{
			BackgroundColor:   r.Styles.BackgroundColor,
			OnBackgroundColor: r.Styles.OnBackgroundColor,
			PrimaryColor:      r.Styles.PrimaryColor,
			OnPrimaryColor:    r.Styles.OnPrimaryColor,
			SecondaryColor:    r.Styles.SecondaryColor,
			OnSecondaryColor:  r.Styles.OnSecondaryColor,
			PrimaryQrCodeData: r.Styles.PrimaryQrCodeData,
		}
		if r.Styles.Rows != nil {
			rows := make([]api.Row, len(r.Styles.Rows))
			for i, row := range r.Styles.Rows {
				rows[i] = api.Row{Left: row.Left, Right: row.Right}
			}
			style.Rows = &rows
		}
		req.Styles = style
	}

	return req
}
