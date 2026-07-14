// Package address parses free-form address strings into structured components
// (street, house number, postcode, city, country).
//
// [ParseAddress] returns the single best interpretation, while
// [ParseAndExpandAddress] returns every plausible interpretation (useful when
// spellings vary).
//
//	res, respStatus, err := address.ParseAddress(client, ctx, sdk.String("Hauptstraße 5A, 12345 Berlin, Germany"))
//	if err == nil && respStatus.HttpStatus == 200 {
//		fmt.Println(res.Street, res.HouseNumber, res.Postcode, res.City)
//	}
//
// See docs/address.md for more examples.
package address
