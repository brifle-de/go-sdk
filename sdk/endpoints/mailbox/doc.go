// Package mailbox searches the documents you have received (inbox) or sent
// (outbox).
//
// Both searches are paginated and accept an optional filter on subject, state
// and type.
//
//	page := float32(1)
//	search := mailbox.InboxSearch{Page: &page}
//	res, respStatus, err := mailbox.SearchMyInbox(client, ctx, &search)
//	if err == nil && respStatus.HttpStatus == 200 {
//		for _, item := range res.Results {
//			fmt.Println(*item.Id, *item.Subject)
//		}
//	}
//
// See docs/mailbox.md for more examples.
package mailbox
