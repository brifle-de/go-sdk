// Package signatures creates signature references and exports signatures.
//
// A signature reference is a reusable definition of who must sign which fields;
// it can then be attached to a document via content.SendContent.
//
//	opts := signatures.SignatureReferenceOptions{
//		Fields: []signatures.SignatureReferenceField{
//			{Name: "signature_1", Purpose: "approval", Role: "customer"},
//		},
//	}
//	res, respStatus, err := signatures.CreateSignatureReference(client, ctx, &tenantID, &opts)
//	if err == nil && respStatus.HttpStatus == 200 {
//		fmt.Println("reference id:", res.Id)
//	}
//
// See docs/signatures.md for more examples.
package signatures
