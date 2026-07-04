// Package content sends and reads documents — the core of the Brifle API.
//
// It also covers receiver checks ([CheckReceiver], [CheckReceiverBulk]),
// delivery certificates ([GetDeliveryCertificate]), delivery status
// ([GetDeliveryStatus]), paper-mail previews ([PreviewPaperMail]) and cover
// letter management ([UploadCoverLetter], [ListCoverLetters], [GetCoverLetter],
// [DeleteCoverLetter]).
//
// # Sending a document
//
//	raw, _ := os.ReadFile("welcome.pdf")
//	req := content.SendContentRequest{
//		To: &content.ReceiverData{
//			BirthInformation: &content.BirthInformationReceiver{
//				FirstName:    sdk.String("Max"),
//				LastName:     sdk.String("Mustermann"),
//				DateOfBirth:  sdk.String("1999-12-12"),
//				PlaceOfBirth: sdk.String("Berlin"),
//			},
//		},
//		Type:    sdk.String(content.Letter),
//		Subject: sdk.String("Welcome to Brifle"),
//		Body: &[]content.ContentItem{
//			{Content: sdk.Base64Encode(raw), Type: sdk.String("application/pdf")},
//		},
//	}
//	res, respStatus, err := content.SendContent(client, ctx, &tenant, &req)
//
// # Identifying a receiver
//
// [ReceiverData] supports three mutually exclusive ways to address a recipient
// — set exactly one of BirthInformation, Email or Phone. Birth information is
// the most precise and is recommended for sensitive content.
//
// # Document types
//
// Use the [Letter], [Invoice] and [Contract] constants for the document Type.
// Only invoices may carry payment info; only non-invoices may request
// signatures.
//
// See docs/content.md and docs/cover-letters.md for more examples.
package content
