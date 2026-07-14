# Cover Letters

Cover letters are templates prepended to documents that are delivered physically (paper mail). Brifle
provides built-in `default` templates, and each tenant can upload its own `custom` templates.

These functions live in the `content` package.

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/content`

## Constants

```go
content.CoverLetterDefault // "default"
content.CoverLetterCustom  // "custom"

content.FormatPdf    // "pdf"
content.FormatBase64 // "base64"
```

## ListCoverLetters

```go
func ListCoverLetters(client *client.BrifleClient, ctx context.Context, tenant *string) (*CoverLettersList, *api.ResponseStatus, error)
```

Lists the cover letters available to a tenant (both built-in and custom).

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
res, respStatus, err := content.ListCoverLetters(client, ctx, &tenant)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

if res.CoverLetters != nil {
	for _, cl := range *res.CoverLetters {
		fmt.Printf("%s (%s)\n", *cl.DisplayName, *cl.Type)
	}
}
```

## UploadCoverLetter

```go
func UploadCoverLetter(client *client.BrifleClient, ctx context.Context, tenant *string, name *string, contentBase64 *string) (*CoverLetter, *api.ResponseStatus, error)
```

Uploads a custom cover letter template. The content must be base64 encoded — use `sdk.Base64Encode`.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"

raw, err := os.ReadFile("cover_letter.pdf")
if err != nil {
	log.Fatal(err)
}

res, respStatus, err := content.UploadCoverLetter(client, ctx, &tenant, sdk.String("welcome-letter"), sdk.Base64Encode(raw))
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("uploaded cover letter:", res.Id, res.DisplayName)
```

## GetCoverLetter

```go
func GetCoverLetter(client *client.BrifleClient, ctx context.Context, tenant *string, coverType *string, fileName *string, format *string) ([]byte, *api.ResponseStatus, error)
```

Downloads the content of a cover letter. `coverType` is `content.CoverLetterDefault` or
`content.CoverLetterCustom`; `format` is `content.FormatPdf` or `content.FormatBase64`. The content is
returned as raw bytes (the PDF itself, or base64 text when using `FormatBase64`).

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"

pdf, respStatus, err := content.GetCoverLetter(
	client, ctx, &tenant,
	sdk.String(content.CoverLetterCustom),
	sdk.String("welcome-letter"),
	sdk.String(content.FormatPdf),
)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

if err := os.WriteFile("downloaded_cover_letter.pdf", pdf, 0o644); err != nil {
	log.Fatal(err)
}
```

## DeleteCoverLetter

```go
func DeleteCoverLetter(client *client.BrifleClient, ctx context.Context, tenant *string, name *string) (*api.ResponseStatus, error)
```

Deletes a custom cover letter template by name.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenant := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
respStatus, err := content.DeleteCoverLetter(client, ctx, &tenant, sdk.String("welcome-letter"))
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}
```
