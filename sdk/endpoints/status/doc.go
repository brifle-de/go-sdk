// Package status reports the health and capabilities of the Brifle API.
//
// This endpoint does not require authentication.
//
//	res, respStatus, err := status.GetStatus(client, ctx)
//	if err != nil {
//		// handle error
//	}
//	if respStatus.HttpStatus == 200 {
//		fmt.Println("service:", *res.Service, "version:", *res.Version)
//	}
//
// See docs/status.md for more examples.
package status
