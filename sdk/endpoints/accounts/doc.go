// Package accounts looks up basic public information about a Brifle account.
//
//	accountID := "2802510314782548"
//	res, respStatus, err := accounts.GetBasicInformation(client, ctx, &accountID)
//	if err == nil && respStatus.HttpStatus == 200 && res.LastName != nil {
//		fmt.Println("last name:", *res.LastName)
//	}
//
// The populated fields depend on the account type (private vs. business).
//
// See docs/accounts.md for more examples.
package accounts
