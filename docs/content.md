# Content

Sending and reading documents is the core of the Brifle API. This package also covers receiver
checks, delivery certificates, delivery status and paper-mail previews.

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/content`

Related: [Cover Letters](cover-letters.md) lives in the same package.

## Document type constants

```go
content.Letter   // "letter"
content.Invoice  // "invoice"
content.Contract // "contract"
```

## SendContent

```go
func SendContent(client *client.BrifleClient, ctx context.Context, tenant *string, req *SendContentRequest) (*SendDocumentResponse, *api.ResponseStatus, error)
```

Sends a document to a receiver on behalf of a tenant.

### Minimal example (send a PDF to a person identified by birth information)

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"

raw, err := os.ReadFile("welcome.pdf")
if err != nil {
	log.Fatal(err)
}

req := content.SendContentRequest{
	To: &content.ReceiverData{
		BirthInformation: &content.BirthInformationReceiver{
			FirstName:    sdk.String("Max"),
			LastName:     sdk.String("Mustermann"),
			DateOfBirth:  sdk.String("1999-12-12"),
			PlaceOfBirth: sdk.String("Berlin"),
		},
	},
	Type:    sdk.String(content.Letter),
	Subject: sdk.String("Welcome to Brifle"),
	Body: &[]content.ContentItem{
		{Content: sdk.Base64Encode(raw), Type: sdk.String("application/pdf")},
	},
}

res, respStatus, err := content.SendContent(client, ctx, &tenant, &req)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("document id:", *res.Id)
```

### Identifying the receiver

`ReceiverData` supports three mutually exclusive ways to address a recipient. Set exactly one:

```go
// 1. By birth information (most precise; good for sensitive content):
content.ReceiverData{
	BirthInformation: &content.BirthInformationReceiver{
		FirstName: sdk.String("Max"), LastName: sdk.String("Mustermann"),
		DateOfBirth: sdk.String("1999-12-12"), PlaceOfBirth: sdk.String("Berlin"),
	},
}

// 2. By email:
content.ReceiverData{
	Email: &content.EmailReceiver{
		Email: sdk.String("max@example.com"),
		Name:  sdk.String("Max Mustermann"),
	},
}

// 3. By phone:
content.ReceiverData{
	Phone: &content.PhoneReceiver{
		PhoneNumber: sdk.String("+491234567890"),
		Name:        sdk.String("Max Mustermann"),
	},
}
```

### Physical delivery fallback (paper mail)

If the receiver cannot be reached electronically, Brifle can fall back to physical delivery. Provide a
postal address:

```go
req.Fallback = &content.Fallback{
	EnabledPhysicalDelivery: true,
	PaperMail: &content.PaperMail{
		Recipient: &content.Recipient{
			AddressLine1: sdk.String("Test Street 1"),
			City:         sdk.String("Berlin"),
			Country:      sdk.String("DE"),
			PostalCode:   sdk.String("12345"),
		},
	},
}
```

### Invoices with payment information

Only documents of type `content.Invoice` may carry payment info:

```go
req.Type = sdk.String(content.Invoice)
req.PaymentInfo = &content.PaymentInfo{
	Payable: boolPtr(true),
	Details: &content.PaymentDetails{
		Amount:    f32Ptr(100), // smallest currency unit (e.g. cents)
		Currency:  sdk.String("EUR"),
		Iban:      sdk.String("DE89370400440532013000"),
		Reference: sdk.String("123456789"),
		DueDate:   sdk.String("2026-12-31"),
	},
}
```

### Requesting signatures

Non-invoice documents may request signatures, either inline or via a
[signature reference](signatures.md):

```go
req.SignatureInfo = &content.SignatureInfo{
	SignatureReference: sdk.String("signature-reference-id"),
}
```

## CheckReceiver

```go
func CheckReceiver(client *client.BrifleClient, ctx context.Context, receiver *ReceiverData) (*ReceiverCheckResponse, *api.ResponseStatus, error)
```

Checks whether a receiver exists on Brifle before sending.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

receiver := content.ReceiverData{
	BirthInformation: &content.BirthInformationReceiver{
		FirstName: sdk.String("Max"), LastName: sdk.String("Mustermann"),
		DateOfBirth: sdk.String("1999-12-12"), PlaceOfBirth: sdk.String("Berlin"),
	},
}

res, respStatus, err := content.CheckReceiver(client, ctx, &receiver)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus == 200 {
	fmt.Println("receiver type:", *res.Receiver.Type)
}
```

## CheckReceiverBulk

```go
func CheckReceiverBulk(client *client.BrifleClient, ctx context.Context, receivers *[]ReceiverData) (*ReceiverBulkCheckResponse, *api.ResponseStatus, error)
```

Checks multiple receivers in one request. Results are returned in the same order as the input.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

receivers := []content.ReceiverData{
	{Email: &content.EmailReceiver{Email: sdk.String("a@example.com"), Name: sdk.String("Alice Example")}},
	{Phone: &content.PhoneReceiver{PhoneNumber: sdk.String("+491234567890"), Name: sdk.String("Bob Example")}},
}

res, respStatus, err := content.CheckReceiverBulk(client, ctx, &receivers)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus == 200 && res.Receivers != nil {
	for i, r := range *res.Receivers {
		if r.Type != nil {
			fmt.Printf("receiver %d found as %s\n", i, *r.Type)
		}
	}
}
```

## GetContent

```go
func GetContent(client *client.BrifleClient, ctx context.Context, documentId *string, readFlag *bool) (*DocumentResponse, *api.ResponseStatus, error)
```

Retrieves a document by ID. Set `readFlag` to `true` to mark it as read.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

documentId := "53C9084932FA27B068424A5FCA81974873E54BC88AAB3B5CCB45C4E6E2C90BB1"
readFlag := false

res, respStatus, err := content.GetContent(client, ctx, &documentId, &readFlag)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

for _, part := range *res.Content {
	fmt.Println("content-type:", *part.ContentType)
	// part.Content is the base64 encoded document
}
```

## GetContentAction

```go
func GetContentAction(client *client.BrifleClient, ctx context.Context, documentId *string) (*ContentActions, *api.ResponseStatus, error)
```

Returns the actions available on a document — pending payments and signature requests.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

documentId := "14288A5F91EA1F2843A2EDEA542E26987F060314D1EF71EB6456CB88865DDA38"
res, respStatus, err := content.GetContentAction(client, ctx, &documentId)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus == 200 && res.Payments != nil && res.Payments.Link != nil {
	fmt.Println("payment link:", *res.Payments.Link)
}
```

## GetDeliveryCertificate

```go
func GetDeliveryCertificate(client *client.BrifleClient, ctx context.Context, documentId *string) (*DeliveryCertificate, *api.ResponseStatus, error)
```

Retrieves the delivery certificate (an Advanced Electronic Seal, XML) for a document.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

documentId := "53C9084932FA27B068424A5FCA81974873E54BC88AAB3B5CCB45C4E6E2C90BB1"
res, respStatus, err := content.GetDeliveryCertificate(client, ctx, &documentId)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus == 200 && res.Certificate != nil {
	fmt.Println(*res.Certificate) // XML
}
```

## GetDeliveryStatus

```go
func GetDeliveryStatus(client *client.BrifleClient, ctx context.Context, documentId *string) (*DeliveryStatus, *api.ResponseStatus, error)
```

Returns the delivery status of a document. The response distinguishes between electronic (`brifle`)
and physical (`physical`) delivery.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

documentId := "53C9084932FA27B068424A5FCA81974873E54BC88AAB3B5CCB45C4E6E2C90BB1"
res, respStatus, err := content.GetDeliveryStatus(client, ctx, &documentId)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

ds := res.DeliveryStatus
fmt.Println("mode:", ds.DeliveryMode)
if ds.Brifle != nil {
	fmt.Println("read:", ds.Brifle.Read, "delivered:", ds.Brifle.DeliveredDate)
}
if ds.Physical != nil {
	fmt.Println("physical state:", ds.Physical.State)
}
```

## PreviewPaperMail

```go
func PreviewPaperMail(client *client.BrifleClient, ctx context.Context, tenant *string, req *PreviewPaperMailRequest) ([]byte, *api.ResponseStatus, error)
```

Renders how a document would look when delivered as physical paper mail and returns the resulting
**PDF as raw bytes**.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"

raw, err := os.ReadFile("welcome.pdf")
if err != nil {
	log.Fatal(err)
}

req := content.PreviewPaperMailRequest{
	To: &content.PreviewReceiver{
		AddressLine1: sdk.String("Test Street 1"),
		City:         sdk.String("Berlin"),
		Country:      sdk.String("DE"),
		PostalCode:   sdk.String("12345"),
	},
	CoverLetter: &content.PreviewCoverLetter{
		Enable: true,
		Name:   sdk.String("default"),
		Type:   sdk.String(content.CoverLetterDefault),
	},
	Body: &content.PreviewBody{
		Content: sdk.Base64Encode(raw),
		Type:    sdk.String("application/pdf"),
	},
}

pdf, respStatus, err := content.PreviewPaperMail(client, ctx, &tenant, &req)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

_ = os.WriteFile("preview.pdf", pdf, 0o644)
```

---

The examples above use two tiny local helpers for pointer values that the `sdk` package does not
provide:

```go
func boolPtr(b bool) *bool       { return &b }
func f32Ptr(f float32) *float32  { return &f }
```
