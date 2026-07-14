// Package wallet issues, reads and revokes wallet items.
//
// A wallet item is a verifiable credential (proof of ownership or permission)
// that can be assigned to a user and optionally exported to Apple or Google
// Wallet.
//
// This is an experimental feature; for more information contact the Brifle team.
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
//	// note: CreateWalletItem returns HTTP 201 on success
//	if err == nil && respStatus.HttpStatus == 201 {
//		fmt.Println("created wallet item:", *res.Id)
//	}
//
// See docs/wallet.md for more examples.
package wallet
