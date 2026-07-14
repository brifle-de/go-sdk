# Signatures

Create signature references (reusable definitions of who must sign which fields) and export existing
signatures.

Import: `github.com/brifle-de/brifle-sdk/sdk/endpoints/signatures`

## CreateSignatureReference

```go
func CreateSignatureReference(client *client.BrifleClient, ctx context.Context, tenantId *string, opts *SignatureReferenceOptions) (*SignatureReference, *api.ResponseStatus, error)
```

Creates a signature reference for a tenant. Each field describes a signature slot: its `Name`, the
optional `Purpose`, and the `Role` of the signer.

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

tenantId := "567e44de-b6b6-4dac-cbce-c5515031f9ea"
opts := signatures.SignatureReferenceOptions{
	Fields: []signatures.SignatureReferenceField{
		{Name: "signature_1", Purpose: "approval", Role: "customer"},
	},
}

res, respStatus, err := signatures.CreateSignatureReference(client, ctx, &tenantId, &opts)
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println("signature reference id:", res.Id)
```

## ExportSignature

```go
func ExportSignature(client *client.BrifleClient, ctx context.Context, signatureId *string, opts *ExportOptions) (*string, *api.ResponseStatus, error)
```

Exports a signature by its ID. The response is returned as a raw string (currently XML). Set the
format via `ExportOptions.Format` (`"xml"`).

```go
ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
defer cancel()

signatureId := "1RTIJVJEAL9BlzdknOoPbgHlg5p0IKRizQaPACh0-O8="
res, respStatus, err := signatures.ExportSignature(client, ctx, &signatureId, &signatures.ExportOptions{
	Format: "xml",
})
if err != nil {
	log.Fatal(err)
}
if respStatus.HttpStatus != 200 {
	log.Fatalf("unexpected status: %d", respStatus.HttpStatus)
}

fmt.Println(*res) // XML document
```
