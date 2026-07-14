// Package tenants lists and fetches the tenants an account owns.
//
// A tenant represents a sending identity (a person or organization) you can
// send content on behalf of. Many endpoints (for example content.SendContent)
// require a tenant ID.
//
//	res, respStatus, err := tenants.GetMyTenants(client, ctx)
//	if err == nil && respStatus.HttpStatus == 200 {
//		for _, t := range res.Tenants {
//			fmt.Println("tenant:", *t.Id)
//		}
//	}
//
// See docs/tenants.md for more examples.
package tenants
